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

func GetPrice(c *gin.Context) {
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

	var price Models.Price

	rows, err := Config.DB.Model(&Models.Price{}).Order("created_at").Rows()
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		Config.DB.ScanRows(rows, &price)

		priceJson, _ := json.Marshal(price)
		err = ws.WriteJSON(string(priceJson))
		if err != nil {
			log.Println("error write json: " + err.Error())
			break
		}
		time.Sleep(1 * time.Second)
	}
}
