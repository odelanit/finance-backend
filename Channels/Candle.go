package Channels

import (
	"../Config"
	"../Models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func GetOHLC(c *gin.Context) {
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		//ReadBufferSize:  1024,
		//WriteBufferSize: 1024,
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error get connection")
		log.Fatal(err)
	}
	defer ws.Close()

	var candle Models.Candle

	rows, err := Config.DB.Model(&Models.Candle{}).Order("time").Rows()
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		Config.DB.ScanRows(rows, &candle)

		candleJson, _ := json.Marshal(candle)
		err = ws.WriteJSON(string(candleJson))
		if err != nil {
			log.Println("error write json: " + err.Error())
			break
		}
		time.Sleep(1 * time.Second)
	}
}
