package port

import (
	"billing-service/internal/user/domain"
	"billing-service/internal/user/event"
	"context"
)

type Service interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByID(ctx context.Context, ID domain.UserID) (domain.User, error)
	// TODO: add lock
	CreditUserBalance(ctx context.Context, ID domain.UserID, amount float64) error
	DebitUserBalance(ctx context.Context, body []byte) (event.SMSUpdateEvent, error)
	UnsuccessfulSMS(ctx context.Context, body []byte) error
	UpdateSMSStatus(ctx context.Context, sms event.SMSUpdateEvent) error
}
