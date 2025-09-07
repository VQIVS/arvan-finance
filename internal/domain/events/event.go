package events

import (
	"context"
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

type RequestSMSBilling struct {
	UserID    string    `json:"user_id"`
	SMSID     string    `json:"sms_id"`
	Amount    int64     `json:"amount"`
	TimeStamp time.Time `json:"timestamp"`
}

type RequestBillingRefund struct {
	TransactionID string    `json:"transaction_id"`
	TimeStamp     time.Time `json:"timestamp"`
}

type SMSDebited struct {
	UserID        string    `json:"user_id"`
	SMSID         string    `json:"sms_id"`
	Amount        int64     `json:"amount"`
	TransactionID string    `json:"transaction_id"`
	TimeStamp     time.Time `json:"timestamp"`
}

func (e *RequestSMSBilling) EventType() EventType {
	return EventTypeDebit
}

func (e *RequestSMSBilling) AggregateID() string {
	return e.UserID
}

func (e *RequestBillingRefund) EventType() EventType {
	return EventTypeRefund
}

func (e *RequestBillingRefund) AggregateID() string {
	return e.TransactionID
}

func (e *SMSDebited) EventType() EventType {
	return EventTypeSMSDebited
}

func (e *SMSDebited) AggregateID() string {
	return e.TransactionID
}
