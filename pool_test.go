package amqp4go

import (
	"fmt"
	"github.com/streadway/amqp"
	"testing"
)

func TestNewAMQP(t *testing.T) {
	var p = NewAMQP("amqp://admin:yangfeng@tw.smartwalle.tk:5672", 2, 1)
	var s = p.GetSession()

	s.h = func(channel *amqp.Channel, d amqp.Delivery) {
		fmt.Println(d.DeliveryTag, string(d.Body), d.ConsumerTag)
		d.Ack(false)
	}
	s.Consume("task_queue", "my_tag", false, false, false, false, nil)

	s = p.GetSession()

	msg := amqp.Publishing{}
	msg.ContentType = "text/plain"
	msg.Body = []byte("exchange")
	s.Channel().Publish("test", "testkey", true, false, msg)

	fmt.Println(s)
	p.Release(s)
}
