package server

import (
	"bufio"
	ut "clientPart/utilities"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	RegAddr   = "http://121.196.155.183:8008/register"
	LoginAddr = "http://121.196.155.183:8008/login"
	WebSocket = "ws://121.196.155.183:8008/index"
	ListAddr  = "http://121.196.155.183:8008/list"
)

var (
	Conn     *websocket.Conn
	UserName string
	Rs       Rooms
)

func LinkServer() error {
	fmt.Printf("正在连接到服务器……")
	conn, _, err := websocket.DefaultDialer.Dial(WebSocket, nil)
	if err != nil {
		fmt.Println("\r连接服务器失败")
		return err
	}

	Conn = conn
	fmt.Printf("\r")
	return nil
}

func GetList() (bool, *Rooms) {
	resp, err := http.Get(ListAddr)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return false, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	err = json.Unmarshal(body, &Rs)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	return true, &Rs
}

func EnterRoom() bool {
	var RoomID int
	fmt.Printf("请输入你想进入的房间号：")
	fmt.Scanln(&RoomID)
	for _, v := range Rs.RS {
		if v.ID == RoomID {
			msg := Message{
				Type:    10,
				Name:    UserName,
				Content: strconv.Itoa(RoomID),
			}

			err := Conn.WriteJSON(&msg)
			if err != nil {
				ut.Error(Conn, err)
				return false
			} else {
				return true
			}
		}
	}
	fmt.Println("该房间号不存在")
	return false
}

func CreateRoom() bool {
	var RoomID int
	fmt.Println("请输入你想创建的房间号")
	fmt.Scanln(&RoomID)
	for _, v := range Rs.RS {
		if v.ID == RoomID {
			fmt.Println("该房间号已存在")
			return false
		}
	}

	msg := Message{
		Type:    20,
		Name:    UserName,
		Content: strconv.Itoa(RoomID),
	}

	err := Conn.WriteJSON(&msg)
	if err != nil {
		ut.Error(Conn, err)
		return false
	} else {
		return true
	}
}

func Send(End *int) {
	var content string
	for {
		input := bufio.NewScanner(os.Stdin)
		if input.Scan() {
			content = input.Text()
		}

		//发送信息模式
		msg := Message{
			Type:    200,
			Name:    UserName,
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			Content: content,
		}

		//由于不能选择文件
		//请手动切换成传输文件模式
		//file, _ := os.Open("./img/mmexport1611376120988.jpg")
		//content, _ = ut.Base64Encoding(file)
		//msg := Message{
		//	Type:    300,
		//	Name:    UserName,
		//	Time:    time.Now().Format("2006-01-02 15:04:05"),
		//	Content: content,
		//}

		err := Conn.WriteJSON(&msg)
		if err != nil {
			*End = 1
			ut.Error(Conn, err)
			return
		}
	}
}

func Receive(End *int) {
	var msg Message
	for {
		err := Conn.ReadJSON(&msg)
		if err != nil {
			*End = 1
			ut.Error(Conn, err)
			return
		}
		switch msg.Type {
		case 100: //广播
			fmt.Println(msg.Name, ":", msg.Content)
		case 200: //存信息
			fmt.Println(msg.Time, msg.Name, ":", msg.Content)
		case 300: //图片
			err = ut.Base64DeEncoding(msg.Content)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
