package model

type CreateItemRequestBody struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
}

type CreateItemResponseBody struct {
	ID int64 `json:"id"`
}
