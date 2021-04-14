package main

import (
	"Homework-4/control"
	"Homework-4/dao"
	"Homework-4/oauth"
	ut "Homework-4/utilities"
)

func main() {
	ut.LogInit()

	err := dao.MySQLInit()
	if err != nil {
		ut.LogError("MySQLInit", err)
		return
	}

	err = oauth.OAuthInit()
	if err != nil {
		ut.LogError("OAuthInit", err)
		return
	}

	control.RoutersEntrance()
}
