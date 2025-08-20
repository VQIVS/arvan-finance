package rabbit

import (
	"billing-service/pkg/constants"
	"billing-service/pkg/logger"
)

// TODO: do we need lock here?
func (r *Rabbit) Consume(queueName string, handler func([]byte) error) error {
	msgs, err := r.Ch.Consume(
		queueName,
		constants.Exchange,
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
			log := logger.GetTracedLogger()

			log.Info("processing message from queue", "queue", queueName)
			if err := handler(d.Body); err != nil {
				log.Error("failed to process message from queue", "queue", queueName, "error", err)
			} else {
				log.Info("successfully processed message from queue", "queue", queueName)
			}
		}
	}()

	return nil
}
