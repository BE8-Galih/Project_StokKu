package entities

import "gorm.io/gorm"

type Item struct {
	*gorm.Model
	Name        string `json:"name" `
	Stocks      int    `json:"stocks"`
	UserID      uint
	HistoryItem []HistoryItem `gorm:"foreignKey:ItemID;references:id"`
}

type HistoryItem struct {
	*gorm.Model
	Name     string `json:"name"`
	ItemName string `json:"itemName"`
	Qty      int    `json:"qty"`
	ItemID   uint   `json:"itemId"`
	UserID   uint   `json:"userId"`
}
