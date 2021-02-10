package main

import (
	"./Config"
	"./Routers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var err error

func main() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	Config.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	router := Routers.SetupRouter()

	_ = router.Run("localhost:8000")
}
