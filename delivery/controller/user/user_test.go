package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"stokku/entities"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var token string

// TEST CREATE USER
func TestCreateUser(t *testing.T) {
	t.Run("Create Success", func(t *testing.T) {
		e := echo.New()
		requestBody, _ := json.Marshal(map[string]interface{}{
			"name":     "Galih",
			"email":    "galih@gmail.com",
			"password": "12345",
		})

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := NewUserControl(&mockResponse{}, validator.New())
		userController.CreateUser()(context)
		type response struct {
			Data  entities.User
			Token string
		}
		var resp response
		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, resp.Data.Name, "Galih")
	})
	t.Run("error Create User", func(t *testing.T) {
		e := echo.New()
		requestBody, _ := json.Marshal(map[string]interface{}{
			"name":     "galih",
			"email":    "galih@gmail.com",
			"password": "12345",
		})

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := NewUserControl(&errMockResponse{}, validator.New())
		userController.CreateUser()(context)

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
	t.Run("error Bind", func(t *testing.T) {
		e := echo.New()
		requestBody := "Data Bind Salah"

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := NewUserControl(&errMockResponse{}, validator.New())
		userController.CreateUser()(context)

		type response struct {
			Code    int
			Message string
		}
		var resp response
		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, "Cannot Bind Data", resp.Message)
		assert.Equal(t, 415, resp.Code)
	})
	t.Run("error Validate User", func(t *testing.T) {
		e := echo.New()
		requestBody, _ := json.Marshal(map[string]interface{}{
			"name":     "galih",
			"email":    "",
			"password": "12345",
		})

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := NewUserControl(&errMockResponse{}, validator.New())
		userController.CreateUser()(context)

		type response struct {
			Code    int
			Message string
		}
		var resp response

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, "Validate Error", resp.Message)
		assert.Equal(t, 406, resp.Code)
	})
}

// TEST LOGIN USER
func TestLogin(t *testing.T) {
	t.Run("Success Login", func(t *testing.T) {
		e := echo.New()
		requestBody, _ := json.Marshal(map[string]interface{}{
			"email":    "galih@gmail.com",
			"password": "Pass123",
		})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/login")

		controller := NewUserControl(&mockResponse{}, validator.New())
		controller.Login()(context)

		type LoginRespon struct {
			Code    int
			Message string
			Status  string
			Data    interface{}
		}

		var resLogin LoginRespon

		json.Unmarshal([]byte(res.Body.Bytes()), &resLogin)
		log.Warn(resLogin.Data)
		assert.Equal(t, 200, resLogin.Code)
		assert.NotNil(t, resLogin.Data)
		data := resLogin.Data.(map[string]interface{})
		token = data["Token"].(string)
		log.Warn(data)
	})

	t.Run("Error Email and Password", func(t *testing.T) {
		e := echo.New()
		requestBody, _ := json.Marshal(map[string]interface{}{
			"email":    "galih@gmail.com",
			"password": "Pass123",
		})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/login")

		controller := NewUserControl(&errMockResponse{}, validator.New())
		controller.Login()(context)

		type LoginRespond struct {
			Code     int
			Messages string
			Status   string
		}
		var resp LoginRespond

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, "Failed", resp.Status)
		assert.Equal(t, 500, resp.Code)
	})
	t.Run("error Bind", func(t *testing.T) {
		e := echo.New()
		requestBody := "Data Bind Salah"

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := NewUserControl(&errMockResponse{}, validator.New())
		userController.Login()(context)

		type response struct {
			Code    int
			Message string
		}
		var resp response
		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, "Cannot Bind Data", resp.Message)
		assert.Equal(t, 415, resp.Code)
	})
	t.Run("error Validate User", func(t *testing.T) {
		e := echo.New()
		requestBody, _ := json.Marshal(map[string]interface{}{
			"name":     "galih",
			"email":    "",
			"password": "12345",
		})

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := NewUserControl(&errMockResponse{}, validator.New())
		userController.Login()(context)

		type response struct {
			Code    int
			Message string
		}
		var resp response

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, "Validate Error", resp.Message)
		assert.Equal(t, 406, resp.Code)
	})
}

func TestGetUserID(t *testing.T) {
	t.Run("Success Get User ID", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/users/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")

		userController := NewUserControl(&mockResponse{}, validator.New())

		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(userController.GetUserID())(context)

		type respondGetUser struct {
			Code    int
			Message string
			Status  string
			data    interface{}
		}
		var respon respondGetUser

		json.Unmarshal([]byte(res.Body.Bytes()), &respon)
		assert.Equal(t, 200, respon.Code)
	})
	t.Run("Error Convert String", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/users/:id")
		context.SetParamNames("id")
		context.SetParamValues("satu")
		userController := NewUserControl(&errMockResponse{}, validator.New())

		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(userController.GetUserID())(context)

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
	t.Run("Error Not Found User", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/users/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")
		userController := NewUserControl(&errMockResponse{}, validator.New())

		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(userController.GetUserID())(context)

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

// TEST UPDATE USER

func TestUpdateUser(t *testing.T) {
	t.Run("Success Update User", func(t *testing.T) {
		e := echo.New()
		requestBody, _ := json.Marshal(map[string]interface{}{
			"name":     "galih",
			"email":    "galih@gmail.com",
			"password": "Pass123",
		})
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/users/:id")
		context.SetParamNames("id")
		context.SetParamValues("2")

		userController := NewUserControl(&mockResponse{}, validator.New())

		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(userController.UpdateUser())(context)

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
		email := data["email"].(string)
		assert.Equal(t, "galih@gmail.com", email)
	})

	t.Run("Error Access Database", func(t *testing.T) {
		e := echo.New()
		requestBody, _ := json.Marshal(map[string]interface{}{
			"email":    "galih@gmail.com",
			"password": "Pass123",
		})
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/users/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")
		userController := NewUserControl(&errMockResponse{}, validator.New())

		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(userController.UpdateUser())(context)

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
		context.SetPath("/users")

		userController := NewUserControl(&errMockResponse{}, validator.New())
		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(userController.UpdateUser())(context)

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
			"email":    "galih@gmail.com",
			"password": "Pass123",
		})
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(string(requestBody)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/users/:id")
		context.SetParamNames("id")
		context.SetParamValues("satu")
		userController := NewUserControl(&errMockResponse{}, validator.New())

		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(userController.UpdateUser())(context)

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

// Test Delete
func TestDeleteID(t *testing.T) {
	t.Run("Success Delete", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
		res := httptest.NewRecorder()
		context := e.NewContext(req, res)
		context.SetPath("/users/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")
		userController := NewUserControl(&mockResponse{}, validator.New())

		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(userController.DeleteUser())(context)

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
		context.SetPath("/users/:id")
		context.SetParamNames("id")
		context.SetParamValues("1")
		userController := NewUserControl(&errMockResponse{}, validator.New())

		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(userController.DeleteUser())(context)

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
		context.SetPath("/users/:id")
		context.SetParamNames("id")
		context.SetParamValues("satu")
		userController := NewUserControl(&errMockResponse{}, validator.New())

		middleware.JWTWithConfig(middleware.JWTConfig{SigningMethod: "HS256", SigningKey: []byte("K3YT0K3N")})(userController.DeleteUser())(context)

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

type mockResponse struct{}

func (m *mockResponse) CreateUser(NewUser entities.User) (entities.User, error) {
	return entities.User{Model: &gorm.Model{ID: uint(1)}, Name: "Galih", Email: "galih@gmail.com", Password: "12345"}, nil
}

func (m *mockResponse) GetUserID(id int) (entities.User, error) {
	return entities.User{Name: "Galih", Email: "galih@gmail.com", Password: "Pass123"}, nil
}

func (m *mockResponse) UpdateUserID(NewData entities.User, id int) error {
	return nil
}

func (m *mockResponse) DeleteUserID(id int) error {
	return nil
}

func (m *mockResponse) Login(email string, password string) (entities.User, error) {
	return entities.User{Model: &gorm.Model{ID: uint(2)}, Email: "galih@gmail.com", Password: "pass123"}, nil
}

// ERROR

type errMockResponse struct{}

func (m *errMockResponse) CreateUser(NewUser entities.User) (entities.User, error) {
	return entities.User{}, errors.New("Cannot Create User")
}

func (m *errMockResponse) GetUserID(id int) (entities.User, error) {
	return entities.User{}, errors.New("Cannot Get Data User")
}

func (m *errMockResponse) UpdateUserID(NewData entities.User, id int) error {
	return errors.New("Cannot Update User")

}

func (m *errMockResponse) DeleteUserID(id int) error {
	return errors.New("Cannot Delete User")
}

func (m *errMockResponse) Login(email string, password string) (entities.User, error) {
	return entities.User{}, errors.New("Cannot Access Database")
}
