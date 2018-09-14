package amqp4go

import (
	"github.com/smartwalle/pool4go"
	"github.com/streadway/amqp"
)

func NewAMQP(addr string, maxActive, maxIdle int) (p *Pool) {
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
	var s = c.(*Session)
	s.tags = make(map[string]struct{})
	return s
}

func (this *Pool) Release(s *Session) error {
	s.Cancel()
	this.Pool.Release(s, false)
	return nil
}
