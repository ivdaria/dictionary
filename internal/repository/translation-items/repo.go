package translationitems

import (
	"context"
	"errors"
	"fmt"

	"dictionary/internal/entity"
	er "dictionary/internal/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	DB *pgxpool.Pool
}

func NewRepo(DB *pgxpool.Pool) *Repo {
	return &Repo{DB: DB}
}

func (r *Repo) GetItemByID(ctx context.Context, id int64) (*entity.TranslationItem, error) {
	const query = `SELECT id, word, translation FROM items WHERE id = $1`
	var mdl model
	if err := r.DB.QueryRow(ctx, query, id).Scan(&mdl); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, er.ErrNotFound
		}
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
		return er.ErrNoRowsAffected
	}
	return nil
}

func (r *Repo) ListItems(ctx context.Context) ([]*entity.TranslationItem, error) {
	const query = `SELECT id, word, translation FROM items ORDER BY word`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("listItems select:  %w", err)
	}
	defer rows.Close()

	var items []*entity.TranslationItem
	for rows.Next() {
		var mdl model
		if err := rows.Scan(&mdl); err != nil {
			return nil, fmt.Errorf("listItems scan row:  %w", err)
		}
		item := mdl.toTranslationItem()
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listItems rows error:  %w", err)
	}

	return items, nil
}
func (r *Repo) DeleteItem(ctx context.Context, id int64) error {
	const query = `DELETE FROM items WHERE id = $1`

	commandTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query to delete item by id: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return er.ErrNoRowsAffected
	}
	return nil
}
