package translationitems

import (
	"context"
	"dictionary/internal/entity"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	DB *pgx.Conn
}

func NewRepo(DB *pgx.Conn) *Repo {
	return &Repo{DB: DB}
}

var ErrNoRowsAffected = errors.New("no rows affected")

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

func (r *Repo) UpdateItem(ctx context.Context, item *entity.TranslationItem) error {
	const query = `UPDATE items SET word = $1, translation = $2 WHERE id = $3`

	mdl := modelFromTranslationItem(item)
	commandTag, err := r.DB.Exec(ctx, query, mdl.Word, mdl.Translation, mdl.ID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}
	return nil
}
