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
	var p = amqp4go.NewAMQP("amqp://admin:yangfeng@tw.smartwalle.tk:5672", "hh", 2, 1)
	var s = p.GetSession()
	s.Handle(h1)
	s.Run("task_queue", false, false, false, false, nil)

	// 如果 session 有 run，则在释放之前，必需要要先 shutdown，否则释放不成功
	//s.Shutdown()
	//p.Release(s)

	done := make(chan bool)
	<-done
}
