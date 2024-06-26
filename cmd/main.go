package main

import (
	"context"
	"dictionary/internal/config"
	"dictionary/internal/entity"
	translationitems "dictionary/internal/repository/translation-items"
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"math/rand"
	"os"
	"time"

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

	for i := 0; i < 1000; i++ {
		itemToCreate := &entity.TranslationItem{
			ID:          0, //Заполнять не надо
			Word:        String(3),
			Translation: String(3),
		}
		id, err := repo.CreateItem(ctx, itemToCreate)
		if err != nil {
			slog.ErrorContext(ctx, "create random item", slog.String("err", err.Error()))
			panic(err)
		}

		item, err := repo.GetItemByID(ctx, id)
		if err != nil {
			slog.ErrorContext(ctx, "get item by id", slog.String("err", err.Error()))
		}
		fmt.Println(item)
	}

}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}
