package model

type CreateItemRequestBody struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
}

type CreateItemResponseBody struct {
	ID int64 `json:"id"`
}

type UpdateItemRequestBody struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
}

type Item struct {
	ID          int64  `json:"id"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
}

type GetItemByIDResponseBody Item

type ListItemsResponseBody struct {
	Items []Item `json:"items"`
}
