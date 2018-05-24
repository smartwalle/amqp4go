package main

import (
	"fmt"
	"github.com/smartwalle/amqp4go/bak"
	"github.com/streadway/amqp"
)

func main() {

	var c = bak.NewSession("amqp://admin:yangfeng@tw.smartwalle.tk:5672", "", "")
	c.Open()

	c.ExchangeDeclare("pubsub", "fanout", true, false, false, false, nil)

	msg := amqp.Publishing{}
	msg.ContentType = "text/plain"
	msg.Body = []byte("exchange pub")
	fmt.Println(c.Publish("pubsub", "pskey", false, false, msg))
}
