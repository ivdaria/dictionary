package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dictionary/internal/config"
	"dictionary/internal/gateway"
	translationitems "dictionary/internal/repository/translation-items"

	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()

	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	var cfg config.Config
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		panic(err)
	}

	dbDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?&pool_max_conns=%d",
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.Database,
		cfg.DBConfig.MaxConnections,
	)

	conn, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	repo := translationitems.NewRepo(conn)

	appServer := gateway.NewAppServer(&cfg, repo)
	go func() {
		if err := appServer.Run(); err != nil {
			slog.ErrorContext(ctx, "app server is closed with error", slog.String("err", err.Error()))
			panic(err)
		}
	}()

	<-ctx.Done()
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	if err := appServer.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "app server shutdown", slog.String("err", err.Error()))
	}
}
