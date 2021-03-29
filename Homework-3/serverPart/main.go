package main

import (
	"log"
	"serverPart/controller"
	"serverPart/dao"
	"serverPart/server"
)

func main() {
	err := dao.MySQLInit()
	if err != nil {
		log.Println("MySQLInit Error", err)
		return
	}

	go server.FilterUpdate()
	go server.BannerUpdate()

	controller.Entrance()
}
