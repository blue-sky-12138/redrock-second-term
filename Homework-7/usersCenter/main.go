// server.go
package main

import (
	"SecondTerm/Homework-7/usersCenter/dao"
	"SecondTerm/Homework-7/usersCenter/serve"
)

//bgacenter start
//begonia -s -c -r ./
func main() {
	err := dao.MySQLInit()
	if err != nil {
		panic(err)
	}
	serve.Entrance()
}
