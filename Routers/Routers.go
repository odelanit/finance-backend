package Routers

import (
	"../Channels"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ws/price", Channels.GetPrice)
	r.GET("/ws/ohlc", Channels.GetOHLC)
	return r
}
