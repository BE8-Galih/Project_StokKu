package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(config AppConfig) *gorm.DB {
	var dsn string
	//Membuat alamat Database
	dsn = fmt.Sprintf("%s:@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		config.Username,
		config.Host,
		config.DBPort,
		config.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}
