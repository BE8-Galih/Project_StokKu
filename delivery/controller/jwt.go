package controller

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// Membuat Token Baru (JWT)
func CreateToken(UserID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["UserID"] = UserID
	claims["expired"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("K3YT0K3N"))
}

// Mengambil Data User ID yang Ada Di Data Token
func ConsumeJWT(c echo.Context) float64 {
	log.Info()
	user := c.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		UserID := claims["UserID"].(float64)
		return UserID
	}
	return 0
}
