package user

import (
	"billing-service/internal/user/domain"
	"billing-service/internal/user/event"
	"billing-service/internal/user/port"
	"billing-service/pkg/adapters/rabbit"
	"billing-service/pkg/constants"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

// TODO: make err handling better
var (
	ErrUserOnCreate        = errors.New("error on creating new user")
	ErrUserNotFound        = errors.New("user not found")
	ErrInsufficientBalance = errors.New("insufficient balance for debit operation")
)

type service struct {
	repo   port.Repo
	rabbit *rabbit.Rabbit
	logger *slog.Logger
	mu     sync.Mutex
}

func NewService(repo port.Repo, rabbit *rabbit.Rabbit) port.Service {
	return &service{
		repo:   repo,
		rabbit: rabbit,
		logger: slog.Default(),
	}
}
func (s *service) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	user, err := s.repo.Create(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *service) GetUserByID(ctx context.Context, ID domain.UserID) (domain.User, error) {
	user, err := s.repo.GetByID(ctx, uint(ID))
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

	user, err := s.repo.GetByID(ctx, uint(ID))
	if err != nil {
		return ErrUserNotFound
	}

	user.Balance += amount

	err = s.repo.UpdateUserBalance(ctx, ID, user.Balance)
	if err != nil {
		return err
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

	user, err := s.repo.GetByID(ctx, msg.UserID)
	if err != nil {
		return event.SMSUpdateEvent{
			Domain: event.SMS,
			SMSID:  msg.SMSID,
			Status: event.StatusFailed,
		}, ErrUserNotFound
	}
	if user.Balance < msg.Amount {
		return event.SMSUpdateEvent{
			Domain: event.SMS,
			SMSID:  msg.SMSID,
			Status: event.StatusFailed,
		}, ErrInsufficientBalance
	}
	err = s.repo.UpdateUserBalance(ctx, domain.UserID(msg.UserID), user.Balance-msg.Amount)
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

	user, err := s.repo.GetByID(ctx, msg.UserID)
	if err != nil {
		return ErrUserNotFound
	}

	if msg.Amount <= 0 {
		return errors.New("invalid SMS amount")
	}

	err = s.repo.UpdateUserBalance(ctx, domain.UserID(msg.UserID), user.Balance+msg.Amount)
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
	s.logger.With("trace_id", uuid.NewString()).Info("sending sms update event", "sms", sms)
	return s.rabbit.Publish(body, constants.KeySMSUpdate)
}
