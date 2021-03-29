package controller

import (
	"github.com/gin-gonic/gin"
	"serverPart/server"
)

func Entrance() {
	r := gin.Default()
	r.POST("/login", server.Login)
	r.POST("/register", server.Register)

	r.GET("/list", server.List)

	r.GET("/index", server.WebsocketConnect)

	r.Run(":8008")
}
