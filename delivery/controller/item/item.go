package item

import (
	"fmt"
	"net/http"
	"stokku/delivery/controller"
	"stokku/delivery/view"
	ItemV "stokku/delivery/view/item"

	"stokku/entities"
	"stokku/repository/item"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type ItemController struct {
	Repo  item.ItemDBControl
	valid *validator.Validate
}

// Membuat Independency Struct
func NewItemControl(ur item.ItemDBControl, val *validator.Validate) *ItemController {
	return &ItemController{
		Repo:  ur,
		valid: val,
	}
}

// Membuat Method Yang Menampilkan Semua Data Item
func (u *ItemController) GetAllItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := u.Repo.GetAllItem()

		if err != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}
		log.Info()
		return c.JSON(http.StatusOK, ItemV.StatusGetAllOk(res))
	}
}

// Membuat Method yang Menampilkan Data Item Berdasarkan ID
func (u *ItemController) GetItemID() echo.HandlerFunc {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, view.ConvertID())
		}
		res, err2 := u.Repo.GetItemID(id)
		if err2 != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}
		log.Info()
		return c.JSON(http.StatusOK, ItemV.StatusGetIdOk(res))
	}
}

//Membuat Method untuk Membuat Data Item Baru
func (u *ItemController) CreateItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		newItem := ItemV.InsertItem{}
		if err := c.Bind(&newItem); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, view.BindData())
		}

		userID := controller.ConsumeJWT(c)

		if err := u.valid.Struct(newItem); err != nil {
			fmt.Println(newItem)
			fmt.Println("Gagal Validasi")
			return c.JSON(http.StatusInternalServerError, view.Validate())
		}

		Itemnew := entities.Item{Name: newItem.Name, Stocks: newItem.Stocks, UserID: uint(int(userID))}
		res, err2 := u.Repo.CreateItem(Itemnew)

		if err2 != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}
		log.Info()
		return c.JSON(http.StatusCreated, ItemV.StatusCreate(res))
	}
}

//Membuat Method Untuk Mengupdate Data Item
func (u *ItemController) UpdateItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		newEdit := ItemV.InsertItem{}
		if err := c.Bind(&newEdit); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, view.BindData())
		}
		idParam := c.Param("id")
		id, err2 := strconv.Atoi(idParam)
		if err2 != nil {
			log.Error(err2)
			return c.JSON(http.StatusInternalServerError, view.ConvertID())

		}
		res := entities.Item{}
		if newEdit.Name != "" {
			res.Name = newEdit.Name
		}
		if newEdit.Stocks != 0 {
			res.Stocks = newEdit.Stocks
		}

		fmt.Println(res)
		ResultUpdate, errUpdate := u.Repo.UpdateItemID(res, id)

		if errUpdate != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}

		log.Info()
		return c.JSON(http.StatusOK, ItemV.StatusUpdate(ResultUpdate))
	}
}

// Membuat Method Untuk Menghapus Data Item
func (u *ItemController) DeleteItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, view.ConvertID())
		}

		errDelete := u.Repo.DeleteItemID(id)
		if errDelete != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}
		log.Info()
		return c.JSON(http.StatusOK, ItemV.StatusDelete())
	}
}

// Membuat Method Untuk Mengedit Data Stok Ketika Terjadi Pembelian Item
func (u ItemController) BuyItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		item := ItemV.TransactionItem{}

		if err := c.Bind(&item); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, view.BindData())
		}

		if err := u.valid.Struct(&item); err != nil {
			return c.JSON(http.StatusInternalServerError, view.Validate())
		}
		fmt.Println(item)

		selectItem, errSelect := u.Repo.SelectItem(item)

		if errSelect != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}

		ID := controller.ConsumeJWT(c)

		res, errBuy := u.Repo.BuyItem(selectItem, item.Qty)
		if errBuy != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}

		NewHistoryItem := entities.HistoryItem{Name: "Pembelian", ItemName: res.Name, Qty: item.Qty, UserID: uint(int(ID)), ItemID: res.ID}

		AddHistory, errAdd := u.Repo.AddHistoryItem(NewHistoryItem)

		if errAdd != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}

		return c.JSON(http.StatusOK, ItemV.StatusCreateHistory(AddHistory))
	}
}

// Membuat Method Untuk Mengedit Data Stok Ketika Terjadi Penjualan Item
func (u ItemController) SellItem() echo.HandlerFunc {
	return func(c echo.Context) error {

		item := ItemV.TransactionItem{}

		if err := c.Bind(&item); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusInternalServerError, view.BindData())
		}

		if err := u.valid.Struct(&item); err != nil {

			return c.JSON(http.StatusInternalServerError, view.Validate())
		}

		selectItem, errSelect := u.Repo.SelectItem(item)

		if errSelect != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}

		ID := controller.ConsumeJWT(c)

		sellItem, errSell := u.Repo.SellItem(selectItem, item.Qty)
		if errSell != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}

		NewHistoryItem := entities.HistoryItem{Name: "Penjualan", ItemName: sellItem.Name, Qty: item.Qty, UserID: uint(int(ID)), ItemID: sellItem.ID}

		AddHistory, errAdd := u.Repo.AddHistoryItem(NewHistoryItem)

		if errAdd != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}
		return c.JSON(http.StatusOK, ItemV.StatusCreateHistory(AddHistory))
	}
}

// Menampilkan History Penjualan Dan Pembelian 1 Minggu Terakhir
func (u *ItemController) History() echo.HandlerFunc {
	return func(c echo.Context) error {
		AllHistory, err := u.Repo.GetAllHistory()

		if err != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}
		log.Info()

		return c.JSON(http.StatusOK, ItemV.StatusGetAllHistory(AllHistory))
	}
}

// Menampilkan History Penjualan 1 Minggu Terakhir
func (u *ItemController) HistorySell() echo.HandlerFunc {
	return func(c echo.Context) error {
		historySell, err := u.Repo.HistorySell("Penjualan")
		if err != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}
		return c.JSON(http.StatusOK, ItemV.StatusGetAllHistory(historySell))
	}
}

//Menampilkan History Pembelian 1 Minggu Terakhir
func (u *ItemController) HistoryBuy() echo.HandlerFunc {
	return func(c echo.Context) error {
		historyBuy, err := u.Repo.HistoryBuy("Pembelian")
		if err != nil {
			log.Warn("Cannot Access Database")
			return c.JSON(http.StatusInternalServerError, view.InternalServerError())
		}
		return c.JSON(http.StatusOK, ItemV.StatusGetAllHistory(historyBuy))
	}
}
