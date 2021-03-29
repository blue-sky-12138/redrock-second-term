package main

import (
	"clientPart/server"
	"fmt"
	"time"
)

var End int //用于判断服务器是否已关闭

func main() {
	menu()

	err := server.LinkServer()
	if err != nil {
		return
	}

	hall()

	go server.Send(&End)
	go server.Receive(&End)

	for {
		time.Sleep(5 * time.Second)
		if End == 1 {
			return
		}
	}
}

func menu() {
	var (
		choice string
	)

	for {
		fmt.Println("**************************************")
		fmt.Println("请选择你想要进行的操作：")
		fmt.Println("1.登录  2.注册  3.退出")
		fmt.Scanln(&choice)
		switch choice {
		case "1":
			ok := server.Login()
			if ok {
				fmt.Println("登陆成功")
				return
			}
		case "2":
			ok := server.Register()
			if ok {
				fmt.Println("注册成功")
				return
			}
		case "3":
			return
		default:
			fmt.Println("输入数据不合法")
		}
	}
}

func hall() {
	var (
		choice string
	)

	for {
		fmt.Println("**************************************")
		fmt.Println("请选择你想要进行的操作：")
		fmt.Println("1.查看房间列表  2.加入房间  3.创建房间")
		fmt.Scanln(&choice)
		switch choice {
		case "1":
			ok, rooms := server.GetList()
			if ok {
				if len((*rooms).RS) == 0 {
					fmt.Println("现在还没有房间，赶快新建房间吧")
				} else {
					for _, v := range (*rooms).RS {
						fmt.Println("房间号:", v.ID, " 人数:", v.Numbers)
					}
				}
			}
		case "2":
			ok := server.EnterRoom()
			if ok {
				fmt.Println("加入房间成功")
				fmt.Println("**************************************")
				return
			}
		case "3":
			ok := server.CreateRoom()
			if ok {
				fmt.Println("创建房间成功")
				fmt.Println("**************************************")
				return
			}
		default:
			fmt.Println("输入数据不合法")
		}
	}
}
