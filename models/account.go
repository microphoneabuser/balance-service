package models

type Account struct {
	Id      int `json:"id"`
	Balance int `json:"balance"`
}

//для маршалинга получаемого json-файла
type AccountInputHandler struct {
	Id          int     `json:"id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type AccountInput struct {
	Id          int    `json:"id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}
