package main

import (
	"Homework-2/control"
	"Homework-2/dao"
	ut "Homework-2/utilities"
)

func main() {
	ut.LogInit()

	err := dao.MySQLDebug()
	if err != nil {
		ut.LogError("MySQLInit", err)
		return
	}

	control.RoutersEntrance()

}
