package main

import (
	"github.com/streadway/amqp"
	"github.com/smartwalle/amqp4go"
)

func main() {
	var p = amqp4go.NewAMQP("amqp://admin:yangfeng@tw.smartwalle.tk:5672", "hh", 10, 1)

	msg := amqp.Publishing{}
	msg.ContentType = "text/plain"
	msg.Body = []byte("exchange")

	var s = p.GetSession()
	s.Channel().Publish("test", "testkey", true, false, msg)
	p.Release(s)
}
