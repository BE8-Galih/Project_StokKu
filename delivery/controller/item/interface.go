package item

type UserInterface interface {
	GetAllItem() error
	GetItemID() error
	CreateItem() error
	UpdateItem() error
	DeleteItem() error
	BuyItem() error
	SellItem() error
	History() error
	HistorySell() error
	HistoryBuy() error
	AddHistory() error
}
