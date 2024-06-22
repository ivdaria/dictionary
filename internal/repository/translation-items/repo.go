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
	row := r.DB.QueryRow(ctx, query, id)

	var mdl model
	if err := row.Scan(&mdl.ID, &mdl.Word, &mdl.Translation); err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	return mdl.toTranslationItem(), nil
}
