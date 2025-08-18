package rabbit

import "billing-service/pkg/logger"

func (r *Rabbit) Consume(queueName string, handler func([]byte) error) error {
	msgs, err := r.Ch.Consume(
		queueName,
		"",
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
			logger.NewLogger().Info("processing message from queue", "queue", queueName)
			if err := handler(d.Body); err != nil {
				logger.NewLogger().Error("failed to process message from queue", "queue", queueName, "error", err)
			} else {
				logger.NewLogger().Info("successfully processed message from queue", "queue", queueName)
			}
		}
	}()

	return nil
}
