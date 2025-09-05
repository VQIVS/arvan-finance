package events

import (
	"context"
	"math/big"
	"time"
)

type EventType string

const (
	EventTypeDebit      EventType = "Debit"
	EventTypeRefund     EventType = "Refund"
	EventTypeSMSDebited EventType = "SMSDebited"
)

type Publisher interface {
	PublishEvent(ctx context.Context, event SMSEvent) error
}

type SMSEvent interface {
	EventType() EventType
	AggregateID() string
}

type DebitUserBalance struct {
	UserID    string    `json:"user_id"`
	SMSID     string    `json:"sms_id"`
	Amount    big.Int   `json:"amount"`
	TimeStamp time.Time `json:"timestamp"`
}

type RefundTransaction struct {
	TransactionID string    `json:"transaction_id"`
	TimeStamp     time.Time `json:"timestamp"`
}

type SMSDebited struct {
	UserID        string    `json:"user_id"`
	SMSID         string    `json:"sms_id"`
	Amount        big.Int   `json:"amount"`
	TransactionID string    `json:"transaction_id"`
	TimeStamp     time.Time `json:"timestamp"`
}

func (e *DebitUserBalance) EventType() EventType {
	return EventTypeDebit
}

func (e *DebitUserBalance) AggregateID() string {
	return e.UserID
}

func (e *RefundTransaction) EventType() EventType {
	return EventTypeRefund
}

func (e *RefundTransaction) AggregateID() string {
	return e.TransactionID
}

func (e *SMSDebited) EventType() EventType {
	return EventTypeSMSDebited
}

func (e *SMSDebited) AggregateID() string {
	return e.TransactionID
}
