package main

import (
	"github.com/smartwalle/amqp4go"
	"github.com/streadway/amqp"
)

func main() {
	var p = amqp4go.NewAMQP("amqp://admin:yangfeng@tw.smartwalle.tk:5672", 2, 1)
	var s = p.GetSession()

	s.Channel().ExchangeDeclare("pubsub", "fanout", true, false, false, false, nil)

	msg := amqp.Publishing{}
	msg.ContentType = "text/plain"
	msg.Body = []byte("exchange pub")
	s.Channel().Publish("pubsub", "pskey", false, false, msg)

	p.Release(s)
}
