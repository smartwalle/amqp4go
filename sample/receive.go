package main

import (
	"fmt"
	"github.com/smartwalle/amqp4go"
	"github.com/streadway/amqp"
)

func h1(c *amqp.Channel, d amqp.Delivery) {
	fmt.Println(d.DeliveryTag, string(d.Body), d.ConsumerTag)
	d.Ack(false)
}

func main() {
	var p = amqp4go.NewAMQP("amqp://admin:yangfeng@tw.smartwalle.tk:5672", 2, 1)
	var s = p.GetSession()
	s.Consume("task_queue", "Tag1", false, false, false, false, nil, h1)

	//p.Release(s)

	done := make(chan bool)
	<-done
}
