package user

import (
	invoiceDomain "billing-service/internal/invoice/domain"

	invoiceRepo "billing-service/internal/invoice/port"
	"billing-service/internal/user/domain"
	"billing-service/internal/user/event"
	"billing-service/internal/user/port"
	userRepo "billing-service/internal/user/port"
	"billing-service/pkg/adapters/rabbit"
	"billing-service/pkg/constants"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
)

// TODO: make err handling better
var (
	ErrUserOnCreate        = errors.New("error on creating new user")
	ErrUserNotFound        = errors.New("user not found")
	ErrInsufficientBalance = errors.New("insufficient balance for debit operation")
	ErrInvalidCreditAmount = errors.New("invalid credit amount")
	ErrInvoiceOnCreate     = errors.New("error on creating invoice")
)

type service struct {
	userRepo    userRepo.Repo
	invoiceRepo invoiceRepo.Repo
	rabbit      *rabbit.Rabbit
	logger      *slog.Logger
	mu          sync.Mutex
}

func NewService(userRepo userRepo.Repo, invoiceRepo invoiceRepo.Repo, rabbit *rabbit.Rabbit) port.Service {
	return &service{
		userRepo:    userRepo,
		invoiceRepo: invoiceRepo,
		rabbit:      rabbit,
		logger:      slog.Default(),
	}
}

func (s *service) CreateUser(ctx context.Context, user domain.User) (domain.APIKey, error) {
	apiKey, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return "", err
	}
	return apiKey, nil
}

func (s *service) GetUserByID(ctx context.Context, ID domain.UserID) (domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, uint(ID))
	if err != nil {
		return domain.User{}, ErrUserNotFound
	}
	return user, nil
}

func (s *service) CreditUserBalance(ctx context.Context, ID domain.UserID, amount float64) error {
	if amount <= 0 {
		return errors.New("invalid credit amount")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user, err := s.userRepo.GetByID(ctx, uint(ID))
	if err != nil {
		return ErrUserNotFound
	}

	user.Balance += amount

	err = s.userRepo.UpdateUserBalance(ctx, ID, user.Balance)
	if err != nil {
		return err
	}
	err = s.invoiceRepo.Create(ctx, &invoiceDomain.Invoice{
		UserID: uint(ID),
		Amount: amount,
		Type:   string(invoiceDomain.InvoiceTypeCredit),
		Status: string(invoiceDomain.InvoiceStatusCompleted),
	})
	if err != nil {
		return ErrInvoiceOnCreate
	}
	return nil
}

func (s *service) DebitUserBalance(ctx context.Context, body []byte) (event.SMSUpdateEvent, error) {
	var msg event.UserBalanceEvent
	if err := json.Unmarshal(body, &msg); err != nil {
		return event.SMSUpdateEvent{}, err
	}
	if msg.Amount <= 0 {
		return event.SMSUpdateEvent{}, errors.New("invalid debit amount")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	user, err := s.userRepo.GetByID(ctx, msg.UserID)
	if err != nil {
		return event.SMSUpdateEvent{}, ErrUserNotFound
	}
	if user.Balance < msg.Amount {
		return event.SMSUpdateEvent{}, ErrInsufficientBalance
	}
	err = s.userRepo.UpdateUserBalance(ctx, domain.UserID(msg.UserID), user.Balance-msg.Amount)
	if err != nil {
		return event.SMSUpdateEvent{}, err
	}

	return event.SMSUpdateEvent{
		Domain: event.SMS,
		SMSID:  msg.SMSID,
		Status: event.StatusSuccess,
	}, nil
}

func (s *service) UnsuccessfulSMS(ctx context.Context, body []byte) error {
	var msg event.UserBalanceEvent
	if err := json.Unmarshal(body, &msg); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	user, err := s.userRepo.GetByID(ctx, msg.UserID)
	if err != nil {
		return ErrUserNotFound
	}

	if msg.Amount <= 0 {
		return errors.New("invalid SMS amount")
	}

	err = s.userRepo.UpdateUserBalance(ctx, domain.UserID(msg.UserID), user.Balance+msg.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateSMSStatus(ctx context.Context, sms event.SMSUpdateEvent) error {
	body, err := json.Marshal(sms)
	if err != nil {
		return err
	}
	s.logger.Info("sending sms update event", "sms", sms)
	return s.rabbit.Publish(body, constants.QueueSMSUpdate)
}
