package item

import (
	"errors"
	"stokku/delivery/view/item"
	"stokku/entities"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type ItemsDB struct {
	Db *gorm.DB
}

func NewDBItem(db *gorm.DB) *ItemsDB {
	return &ItemsDB{
		Db: db,
	}
}

//  Method Untuk Mengambil Semua Item di Database
func (u *ItemsDB) GetAllItem() ([]entities.Item, error) {
	allItems := []entities.Item{}
	if err := u.Db.Find(&allItems).Error; err != nil {
		log.Warn(err)
		return nil, errors.New("Cannot Access Database")
	}
	log.Info()
	return allItems, nil
}

// Method Untuk Mengambil Item Berdasarkan ID di Database
func (u *ItemsDB) GetItemID(id int) (entities.Item, error) {
	ItemID := entities.Item{}
	if err := u.Db.Where("id = ?", id).Find(&ItemID).Error; err != nil {
		log.Warn(err)
		return entities.Item{}, errors.New("Cannot Access Database")
	}
	return ItemID, nil
}

// Method Untuk Menambahkan Data Item Baru di Database
func (u *ItemsDB) CreateItem(NewItem entities.Item) (entities.Item, error) {
	if err := u.Db.Create(&NewItem).Error; err != nil {
		log.Warn(err)
		return entities.Item{}, errors.New("Cannot Access Database")
	}
	log.Info()
	return NewItem, nil
}

// Method Untuk Mengupdate Data Item di Database
func (u *ItemsDB) UpdateItemID(NewData entities.Item, id int) (entities.Item, error) {
	UpdateItem := entities.Item{}
	if err := u.Db.Where("id = ?", id).Updates(entities.Item{Name: NewData.Name, Stocks: NewData.Stocks}).Find(&UpdateItem).Error; err != nil {
		log.Warn(err)
		return UpdateItem, errors.New("Error Updating Data")
	}
	log.Info()
	return UpdateItem, nil
}

// Method Untuk Menghapus Item Berdasarkan ID di Database
func (u *ItemsDB) DeleteItemID(id int) error {
	if err := u.Db.Delete(&entities.Item{}, id).Error; err != nil {
		log.Warn(err)
		return errors.New("Error Access Database")
	}
	log.Info()
	return nil
}

// Method Untuk Mengambil Item Tertentu di Database
func (u *ItemsDB) SelectItem(itemBuy item.TransactionItem) (entities.Item, error) {
	SelectItem := entities.Item{}
	if err := u.Db.Where("name=?", itemBuy.ItemName).Find(&SelectItem).Error; err != nil {
		log.Warn(err)
		return SelectItem, errors.New("Error Access Database")
	}
	return SelectItem, nil
}

// Method Untuk Melakukan Update Data Qty Stok Item di Database Ketika Terjadi Pembelian
func (u *ItemsDB) BuyItem(itemBuy entities.Item, qty int) (entities.Item, error) {

	UpdateItem := entities.Item{}
	if err := u.Db.Where("name = ?", itemBuy.Name).Updates(entities.Item{Name: itemBuy.Name, Stocks: itemBuy.Stocks + qty}).Find(&UpdateItem).Error; err != nil {
		log.Warn(err)
		return UpdateItem, errors.New("Error Updating Data")
	}
	return UpdateItem, nil
}

// Method Untuk Melakukan Update Data Qty Stok Item di Database Ketika Terjadi Penjualan
func (u *ItemsDB) SellItem(itemSell entities.Item, qty int) (entities.Item, error) {

	UpdateItem := entities.Item{}
	if err := u.Db.Where("name = ?", itemSell.Name).Updates(entities.Item{Name: itemSell.Name, Stocks: itemSell.Stocks - qty}).Find(&UpdateItem).Error; err != nil {
		log.Warn(err)
		return UpdateItem, errors.New("Error Updating Data")
	}
	return UpdateItem, nil
}

// Method Untuk Menambahkan Data History Transaksi Ketika Terjadi Pembelian Maupun Penjualan
func (u *ItemsDB) AddHistoryItem(newHistory entities.HistoryItem) (entities.HistoryItem, error) {
	if err := u.Db.Create(&newHistory).Error; err != nil {
		log.Warn(err)
		return entities.HistoryItem{}, errors.New("Error Access Database")
	}
	return newHistory, nil
}

// Method Untuk Mengambil Data History Semua Transaksi 1 Minggu Terakhir Di Database
func (u *ItemsDB) GetAllHistory() ([]entities.HistoryItem, error) {
	HistoryTrasaction := []entities.HistoryItem{}
	if err := u.Db.Where("created_at > DATE_SUB(?, INTERVAL 7 DAY)", time.Now()).Find(&HistoryTrasaction).Error; err != nil {
		log.Warn(err)
		return []entities.HistoryItem{}, errors.New("Error Access Database")
	}
	if len(HistoryTrasaction) == 0 {
		log.Warn("Data Is Empty")
		return nil, errors.New("Data Is Empty")
	}
	log.Info()
	return HistoryTrasaction, nil
}

// Method Untuk Mengambil Data History Penjualan 1 Minggu Terakhir Di Database
func (u *ItemsDB) HistoryBuy(name string) ([]entities.HistoryItem, error) {
	historyBuy := []entities.HistoryItem{}
	if err := u.Db.Where("name = ? AND created_at > DATE_SUB(?, INTERVAL 7 DAY)", name, time.Now()).Find(&historyBuy).Error; err != nil {
		log.Warn(err)
		return []entities.HistoryItem{}, errors.New("Error Access Database")
	}
	return historyBuy, nil
}

// Method Untuk Mengambil Data History Pembelian 1 Minggu Terakhir Di Database
func (u *ItemsDB) HistorySell(name string) ([]entities.HistoryItem, error) {
	historySell := []entities.HistoryItem{}
	if err := u.Db.Where("name = ? AND created_at > DATE_SUB(?, INTERVAL 7 DAY)", name, time.Now()).Find(&historySell).Error; err != nil {
		log.Warn(err)
		return []entities.HistoryItem{}, errors.New("Error Access Database")
	}
	return historySell, nil
}
