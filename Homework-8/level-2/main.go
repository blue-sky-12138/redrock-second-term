package main

import (
	BlueMQ2 "SecondTerm/Homework-8/level-2/BlueMQ"
	"fmt"
	"log"
)

//localhost:10405
func main() {
	go func() {
		r := BlueMQ2.NewServe()
		r.Run()
	}()

	c, err := BlueMQ2.NewClient("localhost:10405")
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Subscribe("RedRock")
	if err != nil {
		log.Println(err)
		return
	}

	err = c.Publish("RedRock", "Hello RedRock!")
	if err != nil {
		log.Println(err)
		return
	}

	msg, _ := c.Require()
	fmt.Println("\n\n\n\n\n", "接受到信息:", msg.Content)
}
