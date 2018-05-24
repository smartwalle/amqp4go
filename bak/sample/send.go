package main

import (
	"fmt"
	"github.com/smartwalle/amqp4go/bak"
	"github.com/streadway/amqp"
)

func main() {

	var c = bak.NewSession("amqp://admin:yangfeng@tw.smartwalle.tk:5672", "task_queue", "tag1")
	c.Open()

	msg := amqp.Publishing{}
	msg.ContentType = "text/plain"
	msg.Body = []byte("exchange")
	fmt.Println(c.Publish("test", "testkey", true, false, msg))
}
