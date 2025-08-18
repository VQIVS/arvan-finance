package rabbit

import (
	"github.com/streadway/amqp"
)

type Rabbit struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func NewRabbit(url string) (*Rabbit, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &Rabbit{Conn: conn, Ch: ch}, nil
}

func (r *Rabbit) Close() {
	if r.Ch != nil {
		_ = r.Ch.Close()
	}
	if r.Conn != nil {
		_ = r.Conn.Close()
	}
}
