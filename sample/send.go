package main

import (
	"github.com/smartwalle/amqp4go"
	"github.com/streadway/amqp"
)

func main() {
	var p = amqp4go.NewAMQP("amqp://admin:yangfeng@tw.smartwalle.tk:5672", 10, 1)

	msg := amqp.Publishing{}
	msg.ContentType = "text/plain"
	msg.Body = []byte("exchange")

	var s = p.GetSession()
	s.Channel().Publish("test", "testkey", true, false, msg)
	p.Release(s)
}
