package BlueMQ

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type Client struct {
	conn      *websocket.Conn
	serve     *Server
	messages  []Message
	isReceive chan bool
	lock      sync.Mutex
}

func NewClient(addr string) (*Client, error) {
	addr = "ws://" + addr + "/link"
	var c Client
	err := c.link(addr)
	if err != nil {
		log.Println("连接服务器失败")
		return nil, err
	}
	c.isReceive = make(chan bool, 1)
	go c.receive()
	return &c, nil
}

func (c *Client) link(addr string) error {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) publish(msg Message, times int) error {
	if times == 3 {
		return errors.New("发送信息异常")
	}

	err := c.conn.WriteJSON(&msg)
	if err != nil {
		return err
	}

	select {
	case <-c.isReceive:
	case <-time.After(5 * time.Second):
		return c.publish(msg, times+1)
	}

	return nil
}

func (c *Client) Publish(topic string, content string) error {
	msg := Message{Code: 101, Topic: topic, Content: content}
	return c.publish(msg, 0)
}

func (c *Client) receive() error {
	var msg Message
	for {
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			return err
		}

		if msg.Code == 100 {
			c.isReceive <- true
		} else {
			if msg.Code != 0 {
				c.conn.WriteJSON(&Message{Code: 100})
				c.lock.Lock()
				c.messages = append(c.messages, msg)
				c.lock.Unlock()
			}
		}
	}
}

func (c *Client) Subscribe(topic string) error {
	msg := Message{Code: 301, Topic: topic}
	err := c.conn.WriteJSON(&msg)
	if err != nil {
		log.Println("订阅异常")
	}
	return err
}

func (c *Client) Unsubscribe(topic string) error {
	msg := Message{Code: 302, Topic: topic}
	err := c.conn.WriteJSON(&msg)
	if err != nil {
		log.Println("取消订阅异常")
	}
	return err
}

func (c *Client) Require() (Message, error) {
	msg := Message{Code: 200}
	err := c.publish(msg, 0)
	if err != nil {
		return Message{}, err
	}

	Done := make(chan bool, 1)
	go func(done chan bool, c *Client) {
		for len(c.messages) == 0 {
		}
		done <- true
	}(Done, c)

	select {
	case <-Done:
		c.lock.Lock()
		msg = c.messages[0]
		c.messages = c.messages[1:]
		c.lock.Unlock()
	case <-time.After(5 * time.Second):
		return Message{}, errors.New("请求超时")
	}

	return msg, nil
}
