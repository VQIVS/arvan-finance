package messaging

import (
	"context"
	"finance/internal/domain/events"
	"finance/pkg/logger"
	"finance/pkg/rabbit"
)

type WalletPublisher struct {
	publisher *rabbit.Publisher
	log       *logger.Logger
}

func NewWalletPublisher(conn *rabbit.RabbitConn, log *logger.Logger) *WalletPublisher {
	return &WalletPublisher{
		publisher: rabbit.NewPublisher(conn),
		log:       log,
	}
}

func (p *WalletPublisher) PublishEvent(ctx context.Context, event events.SMSEvent) error {
	p.log.Info(
		ctx,
		"Publishing event",
		"event_type", event.EventType(),
		"aggregate_id", event.AggregateID())

	return p.publisher.Publish(rabbit.Exchange, rabbit.SMSBilledRouting, event)
}
