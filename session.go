package amqp4go

import (
	"github.com/smartwalle/errors"
	"github.com/streadway/amqp"
	"sync"
)

type handler func(*amqp.Channel, amqp.Delivery)

type Session struct {
	conn *amqp.Connection
	ch   *amqp.Channel

	tag string

	h    handler
	done chan struct{}

	mu        sync.Mutex
	isRunning bool
}

func (this *Session) Close() (err error) {
	this.mu.Lock()

	if this.isRunning {
		this.mu.Unlock()
		this.Shutdown()
		this.mu.Lock()
	}

	defer this.mu.Unlock()
	if this.ch != nil {
		if err = this.ch.Cancel(this.tag, true); err != nil {
			return err
		}
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

func (this *Session) Tag() string {
	return this.tag
}

func (this *Session) Run(queue string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.isRunning {
		return
	}

	if this.ch == nil {
		return errors.New("connection is closed")
	}

	ds, err := this.ch.Consume(queue, this.tag, autoAck, exclusive, noLocal, noWait, args)
	if err != nil {
		return err
	}
	go this.handle(ds)
	this.isRunning = true
	return nil
}

func (this *Session) Shutdown() {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.isRunning == false {
		return
	}

	this.done <- struct{}{}
}

func (this *Session) Handle(h handler) {
	this.h = h
}

func (this *Session) handle(deliveries <-chan amqp.Delivery) {
	for {
		select {
		case d := <-deliveries:

			if this.h != nil {
				this.h(this.ch, d)
			}
		case <-this.done:
			this.isRunning = false
			return
		}
	}
}
