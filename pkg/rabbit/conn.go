package rabbit

import "github.com/streadway/amqp"

type RabbitConn struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func NewRabbitConn(uri string) *RabbitConn {
	conn, err := amqp.Dial(uri)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return &RabbitConn{
		Conn: conn,
		Ch:   ch,
	}
}

// using amq.topic exchnage (no need to declare exchange before)
func (r *RabbitConn) DeclareBindQueue(name, exchange, routing string) error {
	_, err := r.Ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = r.Ch.QueueBind(
		name,
		routing,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
