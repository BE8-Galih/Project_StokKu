package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(config AppConfig) *gorm.DB {
	var dsn string
	//Membuat alamat Database
	dsn = "root:Kuroko25nara@tcp(db-learn.cb8meadbge6r.ap-southeast-1.rds.amazonaws.com:3306)/Stokku?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}
