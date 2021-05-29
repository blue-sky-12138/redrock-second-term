package BlueMQ

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Server struct {
	bro *Broker
}

func NewServe() *Server {
	return &Server{bro: newBroker()}
}

func (s *Server) Run() {
	router := gin.Default()
	router.GET("/link", func(ctx *gin.Context) {
		up := websocket.Upgrader{
			Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		coon, err := up.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println("connectError:", err)
			return
		} else {
			c := newServeClient(coon, s)
			go c.serveClientReceive()
		}
	})
	router.Run(":10405")
}

func newServeClient(conn *websocket.Conn, s *Server) *Client {
	return &Client{
		conn:      conn,
		serve:     s,
		isReceive: make(chan bool, 1),
	}
}

func (c *Client) addMessage(msg Message) {
	c.lock.Lock()
	c.messages = append(c.messages, msg)
	c.lock.Unlock()
}

func (c *Client) serveSendMessage() {
	c.lock.Lock()
	if len(c.messages) == 0 {
		c.sendMsg(Message{201, "NULL", "没有消息"}, 0)
		c.lock.Unlock()
	} else {
		msg := c.messages[0]
		c.messages = c.messages[1:]
		c.lock.Unlock()
		c.sendMsg(msg, 0)
	}
}

func (c *Client) sendMsg(msg Message, times int) {
	if times == 3 {
		c.lock.Lock()
		c.messages = append([]Message{msg}, c.messages...)
		c.lock.Unlock()
		return
	}

	err := c.conn.WriteJSON(&msg)
	if err != nil {
		return
	}
	select {
	case <-c.isReceive:
	case <-time.After(5 * time.Second):
		c.sendMsg(msg, times+1)
	}
}

func (c *Client) serveClientReceive() {
	var msg Message
	for {
		msg = Message{}
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			return
		}

		if msg.Code == 100 {
			c.isReceive <- true
		} else {
			if msg.Code != 0 {
				c.conn.WriteJSON(&Message{Code: 100})
			}
			switch msg.Code {
			case 101:
				c.serve.bro.broadcast(msg)
			case 200:
				c.serveSendMessage()
			case 201:
			case 301:
				c.serve.bro.addSubscribers(msg.Topic, c)
			case 302:
				c.serve.bro.delSubscribers(msg.Topic, c)
			}

		}
	}
}
