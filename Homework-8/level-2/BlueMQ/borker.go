package BlueMQ

import "github.com/gorilla/websocket"

type Broker struct {
	channels map[string]*channel
}

type channel struct { //订阅频道
	subscribers map[*websocket.Conn]*Client //订阅者
}

func newChannel() channel {
	return channel{subscribers: make(map[*websocket.Conn]*Client)}
}

func newBroker() *Broker {
	return &Broker{channels: make(map[string]*channel)}
}

type Message struct {
	Code    int    `json:"code"`
	Topic   string `json:"topic"`
	Content string `json:"content"`
}

func (b *Broker) addSubscribers(topic string, c *Client) {
	val, ok := b.channels[topic]
	if ok {
		val.subscribers[c.conn] = c
	} else {
		channels := newChannel()
		channels.subscribers[c.conn] = c
		b.channels[topic] = &channels
	}
}

func (b *Broker) delSubscribers(topic string, c *Client) {
	delete(b.channels[topic].subscribers, c.conn)
}

func (b *Broker) broadcast(msg Message) {
	msg.Code = 101
	for _, c := range b.channels[msg.Topic].subscribers {
		c.addMessage(msg)
	}
}
