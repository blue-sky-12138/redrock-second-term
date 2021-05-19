package main

import (
	"SecondTerm/Homework-7/videosRouters/control"
	"SecondTerm/Homework-7/videosRouters/dao"
	"SecondTerm/Homework-7/videosRouters/oauth"
	ut "SecondTerm/Homework-7/videosRouters/utilities"
)

func main() {
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
