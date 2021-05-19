package main

import (
	"SecondTerm/Homework-7/fileRouters/control"
	"SecondTerm/Homework-7/fileRouters/oauth"
	ut "SecondTerm/Homework-7/fileRouters/utilities"
)

func main() {
	err := oauth.OAuthInit()
	if err != nil {
		ut.LogError("OAuthInit", err)
		return
	}

	control.RoutersEntrance()
}
