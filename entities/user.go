package entities

import "gorm.io/gorm"

type User struct {
	*gorm.Model
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	Password    string        `json:"password"`
	Item        []Item        `gorm:"foreignKey:UserID;references:id"`
	HistoryItem []HistoryItem `gorm:"foreignKey:UserID;references:id"`
}
