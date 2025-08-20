package rabbit

import "github.com/google/uuid"

func (r *Rabbit) Consume(queueName string, handler func([]byte) error) error {
	msgs, err := r.Ch.Consume(
		queueName,
		"amqp.topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			r.Logger.With("trace_id", uuid.NewString()).Info("processing message from queue", "queue", queueName)
			if err := handler(d.Body); err != nil {
				r.Logger.With("trace_id", uuid.NewString()).Error("failed to process message from queue", "queue", queueName, "error", err)
			} else {
				r.Logger.With("trace_id", uuid.NewString()).Info("successfully processed message from queue", "queue", queueName)
			}
		}
	}()

	return nil
}
