package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port     int16
	DBPort   int16
	Host     string
	Username string
	Password string
	DBName   string
}

// Mengisi Struct AppConfig dengan Data yang ada di local.env
func AppInit() *AppConfig {
	var app *AppConfig
	app = GetConfig()
	if app == nil {
		log.Fatal("Failed To Init Config")
	}
	return app
}

// Mengambil Data Dari Local.Env
func GetConfig() *AppConfig {
	var res AppConfig
	err := godotenv.Load("local.env")

	if err != nil {
		log.Fatal(err)
		return nil
	}
	ports, _ := strconv.Atoi(os.Getenv("PORT"))
	res.Port = int16(ports)
	portDB, _ := strconv.Atoi(os.Getenv("DBPORT"))
	res.DBPort = int16(portDB)
	res.Host = os.Getenv("HOST")
	res.Username = os.Getenv("NAME")
	res.Password = os.Getenv("PASSWORD")
	res.DBName = os.Getenv("DBNAME")
	return &res
}
