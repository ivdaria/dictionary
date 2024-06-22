package translationitems

import "dictionary/internal/entity"

type model struct {
	ID          int64
	Word        string
	Translation string
}

func (m *model) toTranslationItem() *entity.TranslationItem {
	return &entity.TranslationItem{
		ID:          m.ID,
		Word:        m.Word,
		Translation: m.Translation,
	}
}
