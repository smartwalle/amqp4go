package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/smartwalle/amqp4go"
	"github.com/streadway/amqp"
	"os"
)

func identity() string {
	hostname, err := os.Hostname()
	h := sha1.New()
	fmt.Fprint(h, hostname)
	fmt.Fprint(h, err)
	fmt.Fprint(h, os.Getpid())
	return fmt.Sprintf("%x", h.Sum(nil))
}

func h2(c *amqp.Channel, d amqp.Delivery) {
	fmt.Println(d.DeliveryTag, string(d.Body), d.ConsumerTag)
	c.Ack(d.DeliveryTag, false)
}

func main() {
	var c = amqp4go.NewSession("amqp://admin:yangfeng@tw.smartwalle.tk:5672", identity(), "sub")
	c.Handle(h2)
	c.Open()

	c.ExchangeDeclare("pubsub", "fanout", true, false, false, false, nil)

	c.QueueDeclare(identity(), false, true, false, false, nil)
	c.QueueBind(identity(), "pskey", "pubsub", false, nil)

	c.Run(false, false, false, false, nil)

	done := make(chan bool)
	<-done
}
