package main

import (
	"fmt"
	"stokku/config"
	itemC "stokku/delivery/controller/item"
	userC "stokku/delivery/controller/user"
	"stokku/delivery/routes"
	"stokku/repository/item"
	"stokku/repository/user"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func main() {
	// Inisiasi Alamat Database Dan Connect Ke Database
	databaseAddress := config.AppInit()
	Database := config.ConnectDB(*databaseAddress)

	// Database.AutoMigrate(&entities.User{}, &entities.Item{}, &entities.HistoryItem{})

	// Menghubungkan File Repo User dengan File Controller User
	UserRepo := user.NewDBUser(Database)
	UserControl := userC.NewUserControl(UserRepo, validator.New())

	// Menghubungkan File Repo Item dengan File Controller Item
	ItemRepo := item.NewDBItem(Database)
	ItemControl := itemC.NewItemControl(ItemRepo, validator.New())

	// Initiate Echo
	e := echo.New()

	// Memanggil Function Routes
	routes.Path(e, *UserControl, *ItemControl)

	// Menjalankan Program
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", databaseAddress.Port)))
}
