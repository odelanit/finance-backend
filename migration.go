package main

import (
	"./Config"
	"./Models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var err2 error

func main() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	Config.DB, err2 = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err2 != nil {
		log.Fatal(err2)
	}

	Config.DB.AutoMigrate(&Models.Candle{})
	Config.DB.AutoMigrate(&Models.Price{})
}
