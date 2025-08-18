package rabbit

import (
	"billing-service/pkg/logger"
	"fmt"

	"github.com/streadway/amqp"
)

func (r *Rabbit) Publish(Body []byte, Q string) error {

	q, err := r.Ch.QueueDeclare(
		Q,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	pubErr := r.Ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(Body),
		},
	)
	if pubErr != nil {
		return fmt.Errorf("failed to publish message to queue %s: %v, message: %s", q.Name, pubErr, Body)
	}
	logger.NewLogger().Info("published message to queue", "queue", q.Name, "message", string(Body))

	return nil

}
