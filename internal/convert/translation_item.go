package convert

import (
	"dictionary/internal/entity"
	"dictionary/pkg/gateway/model"
)

func ItemFromCreateItemRequestBody(mdl *model.CreateItemRequestBody) *entity.TranslationItem {
	return &entity.TranslationItem{
		ID:          0,
		Word:        mdl.Word,
		Translation: mdl.Translation,
	}
}
