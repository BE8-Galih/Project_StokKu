package routes

import (
	"stokku/delivery/controller/item"
	"stokku/delivery/controller/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Path(e *echo.Echo, userC user.UserController, itemC item.ItemController) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.Logger())

	user := e.Group("/users")
	user.POST("", userC.CreateUser())
	user.POST("/login", userC.Login())
	user.GET("/:id", userC.GetUserID(), middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte("K3YT0K3N")}))
	user.PUT("/:id", userC.UpdateUser(), middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte("K3YT0K3N")}))
	user.DELETE("/:id", userC.DeleteUser(), middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte("K3YT0K3N")}))

	// Login
	e.POST("/login", userC.Login())

	Items := e.Group("/items")
	// Items.GET("/1", uc.GetAllItem, middleware.BasicAuth(Auth.AuthBasic))

	Items.GET("", itemC.GetAllItem(), middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte("K3YT0K3N")}))
	Items.GET("/:id", itemC.GetItemID(), middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte("K3YT0K3N")}))
	Items.POST("", itemC.CreateItem(), middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte("K3YT0K3N")}))
	Items.DELETE("/:id", itemC.DeleteItem(), middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte("K3YT0K3N")}))
	Items.PUT("/:id", itemC.UpdateItem(), middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte("K3YT0K3N")}))

	e.POST("/buy", itemC.BuyItem(), middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte("K3YT0K3N")}))
	e.POST("/sell", itemC.SellItem(), middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte("K3YT0K3N")}))
	e.GET("/history", itemC.History(), middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte("K3YT0K3N")}))
	e.GET("/historySell", itemC.HistorySell(), middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte("K3YT0K3N")}))
	e.GET("/historyBuy", itemC.HistoryBuy(), middleware.JWTWithConfig(middleware.JWTConfig{SigningKey: []byte("K3YT0K3N")}))
}
