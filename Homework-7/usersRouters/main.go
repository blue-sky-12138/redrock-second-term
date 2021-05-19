package main

import (
	"SecondTerm/Homework-7/usersRouters/control"
	"SecondTerm/Homework-7/usersRouters/oauth"
	ut "SecondTerm/Homework-7/usersRouters/utilities"
)

func main() {
	err := oauth.OAuthInit()
	if err != nil {
		ut.LogError("OAuthInit", err)
		return
	}

	control.RoutersEntrance()
}
