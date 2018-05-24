package amqp4go

import (
	"github.com/smartwalle/errors"
	"github.com/smartwalle/pool4go"
	"github.com/streadway/amqp"
)

func NewAMQP(addr, tag string, maxActive, maxIdle int) (p *Pool) {
	var dialFunc = func() (pool4go.Conn, error) {
		conn, err := amqp.Dial(addr)
		if err != nil {
			return nil, err
		}
		ch, err := conn.Channel()
		if err != nil {
			return nil, err
		}

		var s = &Session{}
		s.conn = conn
		s.ch = ch
		s.tag = tag
		s.done = make(chan struct{})
		s.isRunning = false
		return s, nil
	}

	p = &Pool{}
	p.Pool = pool4go.NewPool(dialFunc)
	p.Pool.SetMaxIdleConns(maxIdle)
	p.Pool.SetMaxOpenConns(maxActive)
	return p
}

type Pool struct {
	*pool4go.Pool
}

func (this *Pool) GetSession() *Session {
	var c, err = this.Pool.Get()
	if err != nil {
		return nil
	}
	return c.(*Session)
}

func (this *Pool) Release(s *Session) error {
	if s.isRunning {
		return errors.New("session still running")
	}
	this.Pool.Release(s, false)
	return nil
}
