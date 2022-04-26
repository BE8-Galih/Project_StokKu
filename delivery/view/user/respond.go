package user

import (
	"net/http"
	"stokku/entities"
)

type LoginRespond struct {
	Data  entities.User
	Token string
}

func StatusGetIdOk(data entities.User) map[string]interface{} {
	return map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success Get Data ID",
		"status":  "Success",
		"data":    data,
	}
}

func StatusCreate(data entities.User) map[string]interface{} {
	return map[string]interface{}{
		"code":    http.StatusCreated,
		"message": "Success Create User",
		"status":  "Success",
		"Users":   data,
	}
}

func StatusUpdate(data entities.User) map[string]interface{} {
	return map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Updated",
		"status":  "Success",
		"data":    data,
	}
}

func StatusDelete() map[string]interface{} {
	return map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Deleted",
		"status":  "Success",
	}
}

func StatusLogin(log LoginRespond) map[string]interface{} {
	return map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Login Berhasil",
		"status":  "Success",
		"data":    log,
	}
}
