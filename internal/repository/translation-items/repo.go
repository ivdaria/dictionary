package translationitems

import (
	"context"
	"dictionary/internal/entity"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	DB *pgx.Conn
}

func NewRepo(DB *pgx.Conn) *Repo {
	return &Repo{DB: DB}
}

func (r *Repo) GetItemByID(ctx context.Context, id int64) (*entity.TranslationItem, error) {
	const query = `SELECT id, word, translation FROM items WHERE id = $1`
	var mdl model
	if err := r.DB.QueryRow(ctx, query, id).Scan(&mdl); err != nil {
		return nil, fmt.Errorf("getItemByID scan row:  %w", err)
	}

	return mdl.toTranslationItem(), nil
}

func (r *Repo) CreateItem(ctx context.Context, item *entity.TranslationItem) (int64, error) {
	const query = `INSERT INTO items(word, translation) VALUES ($1,$2) RETURNING id`

	var id int64
	mdl := modelFromTranslationItem(item)
	if err := r.DB.QueryRow(ctx, query, mdl.Word, mdl.Translation).Scan(&id); err != nil {
		return 0, fmt.Errorf("createItem scan row:  %w", err)
	}

	return id, nil
}
