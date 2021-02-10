package main

import (
	"./Config"
	"./Helpers"
	"./Models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type Candles struct {
	Data      []CandleRes
	Timestamp uint
}

type CandleRes struct {
	ID        uint
	Open      string
	High      string
	Low       string
	Close     string
	Period    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

func doEvery(d time.Duration, f func(time2 time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func fetchData(t time.Time) {
	url := "https://api.coincap.io/v2/candles?exchange=poloniex&interval=m1&baseId=bitcoin&quoteId=bitcoin"

	data := Candles{}

	_ = Helpers.GetJson(url, &data)
	candles := data.Data

	for _, candleData := range candles {
		openVal, _ := strconv.ParseFloat(candleData.Open, 64)
		closeVal, _ := strconv.ParseFloat(candleData.Close, 64)
		highVal, _ := strconv.ParseFloat(candleData.High, 64)
		lowVal, _ := strconv.ParseFloat(candleData.Low, 64)
		candle := Models.Candle{Open: openVal, Close: closeVal, High: highVal, Low: lowVal, Time: candleData.Period}
		Config.DB.Create(&candle)
	}
}

var err3 error

func main() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	Config.DB, err3 = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err3 != nil {
		log.Fatal(err3)
	}

	fetchData(time.Now())
	doEvery(60*time.Second, fetchData)
}
