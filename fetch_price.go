package main

import (
	"./Config"
	"./Models"
	"encoding/json"
	"flag"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type PriceRes struct {
	Litecoin string
	Monero   string
}

var err4 error

func main() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	Config.DB, err4 = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err4 != nil {
		log.Fatal(err4)
	}

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u, err := url.Parse("wss://ws.coincap.io/prices?assets=monero,litecoin")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			price := PriceRes{}
			_ = json.Unmarshal(message, &price)
			if price.Litecoin != "" && price.Monero != "" {
				litecoinPrice, _ := strconv.ParseFloat(price.Litecoin, 64)
				moneroPrice, _ := strconv.ParseFloat(price.Monero, 64)
				candle := Models.Price{Litecoin: litecoinPrice, Monero: moneroPrice}
				Config.DB.Create(&candle)
			}
		}
	}()

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Millisecond):
			}
			return
		}
	}
}
