package translationitems

import (
	"dictionary/internal/entity"
	"github.com/jackc/pgx/v5"
)

type model struct {
	ID          int64
	Word        string
	Translation string
}

func (m *model) ScanRow(rows pgx.Rows) error {
	return rows.Scan(&m.ID, &m.Word, &m.Translation)
}

func (m *model) toTranslationItem() *entity.TranslationItem {
	return &entity.TranslationItem{
		ID:          m.ID,
		Word:        m.Word,
		Translation: m.Translation,
	}
}

func modelFromTranslationItem(item *entity.TranslationItem) *model {
	return &model{
		ID:          item.ID,
		Word:        item.Word,
		Translation: item.Translation,
	}
}
