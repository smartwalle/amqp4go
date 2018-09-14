package amqp4go

import (
	"errors"
	"github.com/streadway/amqp"
	"sync"
)

const (
	K_EXCHANGE_KIND_DIRECT  = "direct" // Direct交换机：转发消息到routingKey指定队列（完全匹配，单播）。
	K_EXCHANGE_KIND_FANOUT  = "fanout" // Fanout交换机：转发消息到所有绑定队列（最快，广播）
	K_EXCHNAGE_KIND_TOPIC   = "topic"  // Topic交换机：按规则转发消息（最灵活，组播）
	K_EXCHNAGE_KIND_HEADERS = "headers"
)

type handler func(*amqp.Channel, amqp.Delivery)

type Session struct {
	conn *amqp.Connection
	ch   *amqp.Channel

	tags map[string]struct{}

	h handler

	mu sync.Mutex
}

func (this *Session) Close() (err error) {
	this.Cancel()

	this.mu.Lock()
	defer this.mu.Unlock()

	if this.ch != nil {
		this.ch.Close()
		this.ch = nil
	}
	if this.conn != nil {
		if err = this.conn.Close(); err != nil {
			return err
		}
		this.conn = nil
	}
	return nil
}

func (this *Session) Connection() *amqp.Connection {
	return this.conn
}

func (this *Session) Channel() *amqp.Channel {
	return this.ch
}

func (this *Session) Consume(queue, tag string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (err error) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if _, ok := this.tags[tag]; ok {
		return
	}
	if this.ch == nil {
		return errors.New("connection is closed")
	}
	this.tags[tag] = struct{}{}

	ds, err := this.ch.Consume(queue, tag, autoAck, exclusive, noLocal, noWait, args)
	if err != nil {
		return err
	}
	go this.handle(tag, ds)
	return nil
}

func (this *Session) Cancel() {
	this.mu.Lock()
	defer this.mu.Unlock()

	for tag := range this.tags {
		this.ch.Cancel(tag, true)
		delete(this.tags, tag)
	}
}

func (this *Session) Handle(h handler) {
	this.h = h
}

func (this *Session) handle(tag string, deliveries <-chan amqp.Delivery) {
	for {
		select {
		case d, ok := <-deliveries:
			if ok == false {
				return
			}
			if this.h != nil {
				this.h(this.ch, d)
			}
		}
	}
}
