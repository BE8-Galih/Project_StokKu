package item

import (
	"stokku/delivery/view/item"
	"stokku/entities"
)

type ItemDBControl interface {
	GetAllItem() ([]entities.Item, error)
	GetItemID(id int) (entities.Item, error)
	CreateItem(NewItem entities.Item) (entities.Item, error)
	UpdateItemID(NewData entities.Item, id int) (entities.Item, error)
	DeleteItemID(id int) error
	SelectItem(itemBuy item.TransactionItem) (entities.Item, error)
	BuyItem(itemBuy item.TransactionItem, id float64, qty int) (entities.HistoryItem, error)
	SellItem(itemSell entities.Item, qty int) (entities.Item, error)
	AddHistoryItem(newHistory entities.HistoryItem) (entities.HistoryItem, error)
	GetAllHistory() ([]entities.HistoryItem, error)
	HistorySell(name string) ([]entities.HistoryItem, error)
	HistoryBuy(name string) ([]entities.HistoryItem, error)
}
