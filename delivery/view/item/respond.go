package item

import (
	"net/http"
	"stokku/entities"
)

func StatusGetAllOk(data []entities.Item) map[string]interface{} {
	return map[string]interface{}{
		"code":     http.StatusOK,
		"messages": "Success Get All data",
		"status":   "Success",
		"data":     data,
	}
}

func StatusGetIdOk(data entities.Item) map[string]interface{} {
	return map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Success Get Data ID",
		"status":  "Success",
		"data":    data,
	}
}

func StatusCreate(data entities.Item) map[string]interface{} {
	return map[string]interface{}{
		"code":     http.StatusCreated,
		"messages": "Success Create Item",
		"status":   "Success",
		"data":     data,
	}
}

func StatusUpdate(data entities.Item) map[string]interface{} {
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

func StatusCreateHistory(data entities.HistoryItem) map[string]interface{} {
	return map[string]interface{}{
		"code":    http.StatusCreated,
		"massage": "Add History Success",
		"status":  "succes",
		"data":    data,
	}
}

func StatusGetAllHistory(data []entities.HistoryItem) map[string]interface{} {
	return map[string]interface{}{
		"code":     http.StatusOK,
		"messages": "Success Get All data",
		"status":   "Success",
		"data":     data,
	}
}
