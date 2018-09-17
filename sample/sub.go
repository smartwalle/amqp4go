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
	var p = amqp4go.NewAMQP("amqp://admin:yangfeng@tw.smartwalle.tk:5672", 2, 1)
	var s = p.GetSession()

	s.Channel().ExchangeDeclare("pubsub", "fanout", true, false, false, false, nil)

	s.Channel().QueueDeclare(identity(), false, true, false, false, nil)
	s.Channel().QueueBind(identity(), "pskey", "pubsub", false, nil)

	s.Consume(identity(), "ss", false, false, false, false, nil, h2)

	//s.Cancel()
	//p.Release(s)

	done := make(chan bool)
	<-done
}
