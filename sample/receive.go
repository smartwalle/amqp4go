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

	var c = amqp4go.NewSession("amqp://admin:yangfeng@tw.smartwalle.tk:5672", "task_queue", "tag")
	c.Handle(h1)
	c.Open()
	c.Run(false, false, false, false, nil)

	done := make(chan bool)
	<-done
}
