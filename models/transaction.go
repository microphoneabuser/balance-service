package models

import (
	"time"
)

type Transaction struct {
	Id          int       `json:"id" db:"id"`
	SenderId    int       `json:"sender_id" db:"sender_id"`
	RecipientId int       `json:"recipient_id" db:"recipient_id"`
	Amount      int       `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
	TimeStamp   time.Time `json:"timestamp" db:"timestamp"`
}

//для маршалинга получаемого json-файла
type TransactionInputHandler struct {
	SenderId    int     `json:"sender_id"`
	RecipientId int     `json:"recipient_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type TransactionInputHandlerId struct {
	Id int `json:"id"`
}

type TransactionInput struct {
	SenderId    int    `json:"sender_id"`
	RecipientId int    `json:"recipient_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type OutputParams struct {
	Limit   int
	Offset  int
	SortCol string
	SortDir string
}

type TransactionForOutput struct {
	Id          int       `json:"id"`
	TimeStamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	SenderId    int       `json:"sender_id"`
	RecipientId int       `json:"recipient_id"`
}
