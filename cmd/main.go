package main

import (
	"context"
	"dictionary/internal/config"
	translationitems "dictionary/internal/repository/translation-items"
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	var cfg config.Config
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		panic(err)
	}

	dbDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.Database,
	)
	conn, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	repo := translationitems.NewRepo(conn)

	item, err := repo.GetItemByID(ctx, 1)
	if err != nil {
		slog.ErrorContext(ctx, "get item by id", slog.String("err", err.Error()))
	}
	fmt.Println(item)

}
