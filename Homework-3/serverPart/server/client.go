package server

import (
	"github.com/gorilla/websocket"
	ut "serverPart/utilities"
	"strconv"
)

type Client struct {
	Name    string
	MsgChan chan Message
	Conn    *websocket.Conn
	RoomID  int
}

func NewClient(conn *websocket.Conn) *Client {
	var (
		msg Message
		c   = Client{
			MsgChan: make(chan Message),
			Conn:    conn,
		}
	)
	//获取首个信息，以进行初始化
	err := conn.ReadJSON(&msg)
	ut.Error(conn, err)

	c.Name = msg.Name
	c.RoomID, _ = strconv.Atoi(msg.Content)

	switch msg.Type {
	case 10: //加入房间
		manage.EnterRoom(&c)
	case 20: //创建房间
		manage.CreateRoom(&c)
	}

	return &c
}

func (c *Client) Receive() {
	var msg Message
	for {
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			ut.Error(c.Conn, err)
			manage.QuitRoom(c)
			break
		}
		manage.Rooms[c.RoomID].Broadcast(msg)
	}
}

func (c *Client) Send() {
	for {
		select {
		case msg := <-c.MsgChan:
			err := c.Conn.WriteJSON(&msg)
			if err != nil {
				ut.Error(c.Conn, err)
				manage.QuitRoom(c)
				break
			}
		}

	}
}
