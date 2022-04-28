package item

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"stokku/delivery/controller"
	"stokku/delivery/view/item"
	"stokku/entities"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

var token string

func TestCreateJWT(t *testing.T) {
	t.Run("Create JWT", func(t *testing.T) {
		token, _ = controller.CreateToken(3)
	})
}

func TestGetAllItem(t *testing.T) {
	t.Run("Get All Success", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items")

		itemController := NewItemControl(&mockItemRepository{}, validator.New())

		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.GetAllItem())(context)

		type response struct {
			Code     int
			Messages string
			Status   string
			Data     []entities.Item
		}
		var resp response

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, resp.Data[0].Name, "mangga")
	})
	t.Run("error Get All Item", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items")

		itemController := NewItemControl(&errMockItemRep{}, validator.New())

		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.GetAllItem())(context)

		type response struct {
			Code     int
			Messages string
			Status   string
		}
		var resp response

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, "Failed", resp.Status)
		assert.Equal(t, 500, resp.Code)
	})
}

func TestGetItemID(t *testing.T) {
	t.Run("Success Get Data Item", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		itemController := NewItemControl(&mockItemRepository{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.GetItemID())(context)

		type respondGetitem struct {
			Code    int
			Message string
			Status  string
			Data    interface{}
		}
		var respon respondGetitem

		json.Unmarshal([]byte(res.Body.Bytes()), &respon)
		assert.Equal(t, 200, respon.Code)
	})
	t.Run("Error Convert String", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items/:id")
		context.SetParamNames("id")
		context.SetParamValues("satu")
		itemController := NewItemControl(&errMockItemRep{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.GetItemID())(context)

		type ErrRespon struct {
			Code    int
			Message string
			Status  string
		}

		var Err ErrRespon
		json.Unmarshal([]byte(res.Body.Bytes()), &Err)
		assert.Equal(t, 500, Err.Code)
		assert.Equal(t, "Failed", Err.Status)
	})
	t.Run("Error Not Found item", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")
		itemController := NewItemControl(&errMockItemRep{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.GetItemID())(context)

		type ErrRespon struct {
			Code    int
			Message string
			Status  string
		}

		var Err ErrRespon
		json.Unmarshal([]byte(res.Body.Bytes()), &Err)
		assert.Equal(t, 500, Err.Code)
		assert.Equal(t, "Failed", Err.Status)
	})
}

// func TestCreateItem(t *testing.T) {
// 	t.Run("Create Success", func(t *testing.T) {
// 		e := echo.New()
// 		requestBody, _ := json.Marshal(map[string]interface{}{
// 			"name":   "Apel",
// 			"stocks": 26,
// 		})

// 		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

// 		res := httptest.NewRecorder()
// 		context := e.NewContext(req, res)
// 		context.SetPath("/items")

// 		itemController := NewItemControl(&mockItemRepository{}, validator.New())
// 		itemController.CreateItem()(context)
// 		type response struct {
// 			Code    string
// 			Message string
// 			Status  string
// 			Data    interface{}
// 		}
// 		var resp response
// 		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
// 		assert.Equal(t, 201, resp.Code)
// 		assert.NotNil(t, resp.Data)
// 	})
// 	t.Run("error Create item", func(t *testing.T) {
// 		e := echo.New()
// 		requestBody, _ := json.Marshal(map[string]interface{}{
// 			"name":   "Mangga",
// 			"stocks": 8,
// 		})

// 		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

// 		res := httptest.NewRecorder()
// 		context := e.NewContext(req, res)
// 		context.SetPath("/items")

// 		itemController := NewItemControl(&errMockItemRep{}, validator.New())
// 		itemController.CreateItem()(context)

// 		type response struct {
// 			Code     int
// 			Messages string
// 			Status   string
// 		}
// 		var resp response

// 		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
// 		assert.Equal(t, "Failed", resp.Status)
// 		assert.Equal(t, 500, resp.Code)
// 	})
// 	t.Run("error Bind", func(t *testing.T) {
// 		e := echo.New()
// 		requestBody := "Data Bind Salah"

// 		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

// 		res := httptest.NewRecorder()
// 		context := e.NewContext(req, res)
// 		context.SetPath("/items")

// 		itemController := NewItemControl(&errMockItemRep{}, validator.New())
// 		itemController.CreateItem()(context)

// 		type response struct {
// 			Code    int
// 			Message string
// 		}
// 		var resp response
// 		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
// 		assert.Equal(t, "Cannot Bind Data", resp.Message)
// 		assert.Equal(t, 415, resp.Code)
// 	})
// 	t.Run("error Validate item", func(t *testing.T) {
// 		e := echo.New()
// 		requestBody, _ := json.Marshal(map[string]interface{}{
// 			"name":     "galih",
// 			"email":    "",
// 			"password": "12345",
// 		})

// 		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

// 		res := httptest.NewRecorder()
// 		context := e.NewContext(req, res)
// 		context.SetPath("/items")

// 		itemController := NewItemControl(&errMockItemRep{}, validator.New())
// 		itemController.CreateItem()(context)

// 		type response struct {
// 			Code    int
// 			Message string
// 		}
// 		var resp response

// 		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
// 		assert.Equal(t, "Validate Error", resp.Message)
// 		assert.Equal(t, 406, resp.Code)
// 	})
// }

func TestUpdateItem(t *testing.T) {
	t.Run("Success Update item", func(t *testing.T) {
		e := echo.New()
		requestBody, _ := json.Marshal(map[string]interface{}{
			"name":   "semangka",
			"stocks": 20,
		})
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items/:id")
		context.SetParamNames("id")
		context.SetParamValues("2")

		itemController := NewItemControl(&mockItemRepository{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.UpdateItem())(context)
		type UpdateRespon struct {
			Code    int
			Message string
			Status  string
			Data    interface{}
		}

		var respon UpdateRespon

		json.Unmarshal([]byte(res.Body.Bytes()), &respon)

		assert.Equal(t, 200, respon.Code)
		assert.Equal(t, "Updated", respon.Message)
		assert.Equal(t, "Success", respon.Status)
		data := respon.Data.(map[string]interface{})
		buah := data["name"].(string)
		assert.Equal(t, "Apel", buah)
	})

	t.Run("Error Access Database", func(t *testing.T) {
		e := echo.New()
		requestBody, _ := json.Marshal(map[string]interface{}{
			"name":   "manggis",
			"stocks": 15,
		})
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")
		itemController := NewItemControl(&errMockItemRep{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.UpdateItem())(context)

		type errUpdate struct {
			Code     int
			Messages string
			Status   string
		}
		var resp errUpdate

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, "Failed", resp.Status)
		assert.Equal(t, 500, resp.Code)
	})
	t.Run("error Bind", func(t *testing.T) {
		e := echo.New()
		requestBody := "Data Bind Salah"

		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)

		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items")

		itemController := NewItemControl(&errMockItemRep{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.UpdateItem())(context)
		type response struct {
			Code    int
			Message string
		}
		var resp response
		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, "Cannot Bind Data", resp.Message)
		assert.Equal(t, 415, resp.Code)
	})
	t.Run("Error Convert String", func(t *testing.T) {
		e := echo.New()

		requestBody, _ := json.Marshal(map[string]interface{}{
			"name":   "semangka",
			"stocks": 20,
		})
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)

		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items/:id")
		context.SetParamNames("id")
		context.SetParamValues("satu")
		itemController := NewItemControl(&errMockItemRep{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.UpdateItem())(context)

		type ErrRespon struct {
			Code    int
			Message string
			Status  string
		}

		var Err ErrRespon
		json.Unmarshal([]byte(res.Body.Bytes()), &Err)
		assert.Equal(t, 500, Err.Code)
		assert.Equal(t, "Failed", Err.Status)
		assert.Equal(t, "Cannot Convert ID", Err.Message)
	})
}

func TestDeleteItem(t *testing.T) {
	t.Run("Success Delete", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")
		itemController := NewItemControl(&mockItemRepository{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.DeleteItem())(context)

		type responDelete struct {
			Code    int
			Message string
			Status  string
		}

		var Err responDelete
		json.Unmarshal([]byte(res.Body.Bytes()), &Err)
		assert.Equal(t, 200, Err.Code)
		assert.Equal(t, "Success", Err.Status)
		assert.Equal(t, "Deleted", Err.Message)
	})
	t.Run("Error Access Database", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")
		itemController := NewItemControl(&errMockItemRep{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.DeleteItem())(context)

		type response struct {
			Code     int
			Messages string
			Status   string
		}
		var resp response

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, "Failed", resp.Status)
		assert.Equal(t, 500, resp.Code)
	})
	t.Run("Error Convert String", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items/:id")
		context.SetParamNames("id")
		context.SetParamValues("satu")
		itemController := NewItemControl(&errMockItemRep{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.DeleteItem())(context)

		type ErrRespon struct {
			Code    int
			Message string
			Status  string
		}

		var Err ErrRespon
		json.Unmarshal([]byte(res.Body.Bytes()), &Err)
		assert.Equal(t, 500, Err.Code)
		assert.Equal(t, "Failed", Err.Status)
		assert.Equal(t, "Cannot Convert ID", Err.Message)
	})
}

func TestSellItem(t *testing.T) {

}

func TestBuyItem(t *testing.T) {

}

func TestHistory(t *testing.T) {
	t.Run("Success Get History", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items")
		itemController := NewItemControl(&mockItemRepository{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.History())(context)
		type responHistory struct {
			Code    int
			Message string
			Status  string
			Data    interface{}
		}

		var respon responHistory
		json.Unmarshal([]byte(res.Body.Bytes()), &respon)
		assert.Equal(t, 200, respon.Code)
		assert.Equal(t, "Success", respon.Status)
		assert.NotNil(t, respon.Data)
	})
	t.Run("Error Get History", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items")
		itemController := NewItemControl(&errMockItemRep{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.History())(context)

		type response struct {
			Code     int
			Messages string
			Status   string
		}
		var resp response

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, "Failed", resp.Status)
		assert.Equal(t, 500, resp.Code)
	})
}

func TestHistorySell(t *testing.T) {
	t.Run("Success Get Sell History", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items")
		itemController := NewItemControl(&mockItemRepository{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.HistorySell())(context)

		type responHistory struct {
			Code    int
			Message string
			Status  string
			Data    interface{}
		}

		var respon responHistory
		json.Unmarshal([]byte(res.Body.Bytes()), &respon)
		assert.Equal(t, 200, respon.Code)
		assert.Equal(t, "Success", respon.Status)
		assert.NotNil(t, respon.Data)
	})
	t.Run("Error Get Sell", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items")
		itemController := NewItemControl(&errMockItemRep{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.HistorySell())(context)

		type response struct {
			Code     int
			Messages string
			Status   string
		}
		var resp response

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, "Failed", resp.Status)
		assert.Equal(t, 500, resp.Code)
	})
}

func TestHistoryBuy(t *testing.T) {
	t.Run("Success Get Buy History", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items")
		itemController := NewItemControl(&mockItemRepository{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.HistoryBuy())(context)
		type responHistory struct {
			Code    int
			Message string
			Status  string
			Data    interface{}
		}

		var respon responHistory
		json.Unmarshal([]byte(res.Body.Bytes()), &respon)
		assert.Equal(t, 200, respon.Code)
		assert.Equal(t, "Success", respon.Status)
		assert.NotNil(t, respon.Data)
	})
	t.Run("Error Get Buy", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/items")
		itemController := NewItemControl(&errMockItemRep{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(itemController.HistoryBuy())(context)
		type response struct {
			Code     int
			Messages string
			Status   string
		}
		var resp response

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, "Failed", resp.Status)
		assert.Equal(t, 500, resp.Code)
	})
}

// Dummy

type mockItemRepository struct{}

func (m *mockItemRepository) GetAllItem() ([]entities.Item, error) {
	return []entities.Item{{Name: "mangga", Stocks: 10}}, nil
}

func (m *mockItemRepository) GetItemID(id int) (entities.Item, error) {
	return entities.Item{Name: "apel", Stocks: 10}, nil
}

func (m *mockItemRepository) CreateItem(NewItem entities.Item) (entities.Item, error) {
	return entities.Item{Name: "Apel", Stocks: 4}, nil
}

func (m *mockItemRepository) UpdateItemID(NewItem entities.Item, id int) (entities.Item, error) {
	return entities.Item{Name: "Apel", Stocks: 20}, nil
}
func (m *mockItemRepository) DeleteItemID(id int) error {
	return nil

}

func (m *mockItemRepository) SelectItem(item item.TransactionItem) (entities.Item, error) {
	return entities.Item{}, nil

}

func (m *mockItemRepository) BuyItem(itemBuy item.TransactionItem, id float64, qty int) (entities.HistoryItem, error) {
	return entities.HistoryItem{}, nil

}

func (m *mockItemRepository) SellItem(itemSell entities.Item, qty int) (entities.Item, error) {
	return entities.Item{}, nil

}

func (m *mockItemRepository) AddHistoryItem(newHistory entities.HistoryItem) (entities.HistoryItem, error) {
	return entities.HistoryItem{}, nil

}

func (m *mockItemRepository) GetAllHistory() ([]entities.HistoryItem, error) {
	return []entities.HistoryItem{{Name: "Pembelian", ItemName: "Baju", Qty: 8}}, nil

}

func (m *mockItemRepository) HistorySell(name string) ([]entities.HistoryItem, error) {
	return []entities.HistoryItem{{Name: "Pembelian", ItemName: "Baju", Qty: 8}}, nil
}

func (m *mockItemRepository) HistoryBuy(name string) ([]entities.HistoryItem, error) {
	return []entities.HistoryItem{{Name: "Penjualan", ItemName: "Celana", Qty: 8}}, nil

}

// Error
type errMockItemRep struct{}

func (m *errMockItemRep) GetAllItem() ([]entities.Item, error) {
	return nil, errors.New("Cannot Access Database")
}

func (m *errMockItemRep) GetItemID(id int) (entities.Item, error) {
	return entities.Item{}, errors.New("Cannot Access Database")
}

func (m *errMockItemRep) CreateItem(NewItem entities.Item) (entities.Item, error) {
	return entities.Item{}, errors.New("Internal Server Error")
}

func (m *errMockItemRep) UpdateItemID(NewItem entities.Item, id int) (entities.Item, error) {
	return entities.Item{}, errors.New("Internal Server Error")
}
func (m *errMockItemRep) DeleteItemID(id int) error {
	return errors.New("Internal Server Error")

}

func (m *errMockItemRep) SelectItem(item item.TransactionItem) (entities.Item, error) {
	return entities.Item{}, nil

}

func (m *errMockItemRep) BuyItem(itemBuy item.TransactionItem, id float64, qty int) (entities.HistoryItem, error) {
	return entities.HistoryItem{}, nil

}

func (m *errMockItemRep) SellItem(itemSell entities.Item, qty int) (entities.Item, error) {
	return entities.Item{}, nil

}

func (m *errMockItemRep) AddHistoryItem(newHistory entities.HistoryItem) (entities.HistoryItem, error) {
	return entities.HistoryItem{}, nil

}

func (m *errMockItemRep) GetAllHistory() ([]entities.HistoryItem, error) {
	return []entities.HistoryItem{}, errors.New("Error Access Database")

}

func (m *errMockItemRep) HistorySell(name string) ([]entities.HistoryItem, error) {
	return []entities.HistoryItem{}, errors.New("Error Access Database")

}

func (m *errMockItemRep) HistoryBuy(name string) ([]entities.HistoryItem, error) {
	return []entities.HistoryItem{}, errors.New("Error Access Database")
}
