package messaging

import (
	"context"
	"finance/internal/domain/events"
	"finance/pkg/rabbit"
)

const (
	exchange             = "amq.topic"
	routingKeySMSDebited = "finance.sms_debited"
)

type WalletPublisher struct {
	publisher *rabbit.Publisher
}

func NewWalletPublisher(conn *rabbit.RabbitConn) *WalletPublisher {
	return &WalletPublisher{
		publisher: rabbit.NewPublisher(conn),
	}
}

func (p *WalletPublisher) PublishEvent(ctx context.Context, event events.SMSEvent) error {
	return p.publisher.Publish(exchange, routingKeySMSDebited, event)
}
