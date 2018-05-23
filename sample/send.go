package main

import (
	"fmt"
	"github.com/smartwalle/amqp4go"
	"github.com/streadway/amqp"
)

func main() {

	var c = amqp4go.NewSession("amqp://admin:yangfeng@tw.smartwalle.tk:5672", "task_queue", "tag1")
	c.Open()

	msg := amqp.Publishing{}
	msg.ContentType = "text/plain"
	msg.Body = []byte("exchange")
	fmt.Println(c.PublishWithExchange("test", "testkey", true, false, msg))

	msg.Body = []byte("queue")
	fmt.Println(c.Publish(true, false, msg))
}
