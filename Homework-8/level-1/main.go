package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

var MQ *amqp.Connection

//这里消息队列产品使用RabbitMQ
//框架采用github.com/streadway/amqp
func main() {
	topic := "RedRock"
	body := "Hello RedRock!"

	channel, err := QueueDeclare(topic)
	if err != nil {
		log.Println(err)
		return
	}

	err = Publish(channel, topic, body)
	if err != nil {
		log.Println(err)
		return
	}

	err = Consume(channel, topic)
	if err != nil {
		log.Println(err)
		return
	}
}

func init() {
	var err error
	MQ, err = amqp.Dial("amqp://blue_sky:135246@localhost:5672/")
	if err != nil {
		return
	}
}

func QueueDeclare(topic string) (*amqp.Channel, error) {
	channel, err := MQ.Channel()
	if err != nil {
		return nil, err
	}

	_, err = channel.QueueDeclare(topic, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func Publish(channel *amqp.Channel, topic string, content string) error {
	body := []byte(content)
	err := channel.Publish("", topic, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})
	return err
}

func Consume(channel *amqp.Channel, topic string) error {
	msgChannel, err := channel.Consume(topic, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	msg := <-msgChannel
	fmt.Printf("接收到信息:%s\n", msg.Body)
	return nil
}
