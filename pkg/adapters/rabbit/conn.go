package rabbit

import (
	"billing-service/pkg/constants"
	"billing-service/pkg/logger"
	"log/slog"

	"github.com/streadway/amqp"
)

type Rabbit struct {
	Conn   *amqp.Connection
	Ch     *amqp.Channel
	Logger *slog.Logger
}

func NewRabbit(url string, customLogger *slog.Logger) (*Rabbit, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	log := customLogger
	if log == nil {
		log = logger.GetLogger()
	}

	return &Rabbit{Conn: conn, Ch: ch, Logger: log}, nil
}

func (r *Rabbit) Close() {
	if r.Ch != nil {
		_ = r.Ch.Close()
	}
	if r.Conn != nil {
		_ = r.Conn.Close()
	}
}

func (r *Rabbit) InitQueues(queue string) error {
	if r == nil || r.Ch == nil {
		return nil
	}
	queueName := GetQueueName(constants.KeyBalanceUpdate)
	_, err := r.Ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = r.Ch.QueueBind(queueName, constants.KeyBalanceUpdate, constants.Exchange, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func GetQueueName(key string) string {
	return "finance_" + key
}
