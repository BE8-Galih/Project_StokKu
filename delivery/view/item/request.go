package item

import "time"

type InsertItem struct {
	Name   string `json:"name" validate:"required"`
	Stocks int    `json:"stocks" validate:"required"`
}

type TransactionItem struct {
	ItemName string `json:"itemName"`
	Qty      int    `json:"qty"`
}

type HistoryTransaction struct {
	Name            string    `json:"name"`
	ItemName        string    `json:"itemName"`
	Qty             int       `json:"qty"`
	TransactionDate time.Time `json:"transactionDate"`
}
