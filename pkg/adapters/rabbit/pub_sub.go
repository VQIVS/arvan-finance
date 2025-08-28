package rabbit

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/streadway/amqp"
)

func (r *Rabbit) Publish(routingKey, exchange string, body []byte) error {
	if r.Ch == nil {
		return errors.New("no channel")
	}

	err := r.Ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish %s: %w", routingKey, err)
	}
	return nil
}

func (r *Rabbit) Subscribe(queueName string, handler func(amqp.Delivery) error) (<-chan amqp.Delivery, error) {
	err := r.Ch.Qos(1, 0, false)
	if err != nil {
		return nil, err
	}

	deliveries, err := r.Ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	go func() {
		for d := range deliveries {
			retryOperation := func() (struct{}, error) {
				return struct{}{}, handler(d)
			}

			backOff := backoff.NewExponentialBackOff()
			backOff.InitialInterval = 1 * time.Second
			backOff.MaxInterval = 30 * time.Second

			_, err := backoff.Retry(context.TODO(), retryOperation, backoff.WithBackOff(backOff))
			if err != nil {
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	return deliveries, nil
}
