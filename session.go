package amqp4go

import (
	"github.com/smartwalle/errors"
	"github.com/streadway/amqp"
	"sync"
)

type handler func(*amqp.Channel, amqp.Delivery)

type Session struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	addr  string
	tag   string
	queue string

	h    handler
	done chan error

	mu        sync.Mutex
	isRunning bool
}

func NewSession(addr, queue, tag string) (s *Session) {
	s = &Session{}
	s.addr = addr
	s.queue = queue
	s.tag = tag
	s.done = make(chan error)
	s.isRunning = false
	return s
}

func (this *Session) Conn() *amqp.Connection {
	return this.conn
}

func (this *Session) Channel() *amqp.Channel {
	return this.channel
}

func (this *Session) Open() (err error) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.isRunning {
		return nil
	}

	defer func() {
		if err != nil && this.channel != nil {
			this.channel.Close()
		}
		if err != nil && this.conn != nil {
			this.conn.Close()
		}
	}()

	if this.conn != nil || this.channel != nil {
		this.mu.Unlock()
		this.Shutdown()
		this.mu.Lock()
	}

	if this.conn, err = amqp.Dial(this.addr); err != nil {
		return err
	}

	if this.channel, err = this.conn.Channel(); err != nil {
		return err
	}
	return nil
}

func (this *Session) Shutdown() (err error) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.isRunning == false {
		return nil
	}

	if this.channel != nil {
		if err = this.channel.Cancel(this.tag, true); err != nil {
			return err
		}
		this.channel = nil
	}

	if this.conn != nil {
		if err = this.conn.Close(); err != nil {
			return err
		}
		this.conn = nil
	}
	return <-this.done
}

func (this *Session) Handle(h handler) {
	this.h = h
}

func (this *Session) Run(autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.isRunning {
		return
	}

	if this.channel == nil {
		return errors.New("connection is closed")
	}

	ds, err := this.channel.Consume(this.queue, this.tag, autoAck, exclusive, noLocal, noWait, args)
	if err != nil {
		return err
	}
	go this.handle(ds)
	this.isRunning = true
	return nil
}

func (this *Session) ExchangeDeclare(exchange, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return this.channel.ExchangeDeclare(exchange, kind, durable, autoDelete, internal, noWait, args)
}

func (this *Session) QueueDeclare(queue string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return this.channel.QueueDeclare(queue, durable, autoDelete, exclusive, noWait, args)
}

func (this *Session) QueueBind(queue, key, exchange string, noWait bool, args amqp.Table) error {
	return this.channel.QueueBind(queue, key, exchange, noWait, args)
}

func (this *Session) Publish(mandatory, immediate bool, msg amqp.Publishing) error {
	return this.channel.Publish("", this.queue, mandatory, immediate, msg)
}

func (this *Session) PublishWithQueue(queue string, mandatory, immediate bool, msg amqp.Publishing) error {
	return this.channel.Publish("", queue, mandatory, immediate, msg)
}

func (this *Session) PublishWithExchange(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return this.channel.Publish(exchange, key, mandatory, immediate, msg)
}

func (this *Session) handle(deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		if this.h != nil {
			this.h(this.channel, d)
		}
	}
	this.done <- nil
	this.isRunning = false
}
