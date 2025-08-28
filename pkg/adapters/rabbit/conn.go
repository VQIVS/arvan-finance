package rabbit

import (
	"billing-service/pkg/constants"
	"billing-service/pkg/logger"
	"fmt"
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
		_ = conn.Close()
		return nil, err
	}

	log := customLogger
	if log == nil {
		log = logger.GetLogger()
	}

	return &Rabbit{Conn: conn, Ch: ch, Logger: log}, nil
}

func (r *Rabbit) Close() error {
	if r.Ch != nil {
		if err := r.Ch.Close(); err != nil {
			return err
		}
	}
	if r.Conn != nil {
		return r.Conn.Close()
	}
	return nil
}

func (r *Rabbit) InitQueues(keys []string) error {
	if r == nil || r.Ch == nil {
		return nil
	}

	for _, queue := range keys {
		queueName := GetQueueName(queue)
		_, err := r.Ch.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}

		if err := r.Ch.QueueBind(queueName, constants.KeySMSUpdate, constants.TopicExchange, false, nil); err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", queueName, err)
		}
	}

	return nil
}

func GetQueueName(key string) string {
	return constants.ServiceName + "_" + key
}
