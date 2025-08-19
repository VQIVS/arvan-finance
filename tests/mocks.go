package tests

import (
	"billing-service/internal/user/domain"
	"billing-service/internal/user/event"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
)

type MockUserRepo struct {
	users          map[uint]domain.User
	mu             sync.Mutex
	createError    error
	getByIDError   error
	updateError    error
	createAPIKeyFn func() domain.APIKey
	nextUserID     uint
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		users:          make(map[uint]domain.User),
		nextUserID:     1,
		createAPIKeyFn: func() domain.APIKey { return "test-api-key" },
	}
}

func (m *MockUserRepo) Create(ctx context.Context, user domain.User) (domain.APIKey, error) {
	if m.createError != nil {
		return "", m.createError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if user.ID == 0 {
		user.ID = domain.UserID(m.nextUserID)
		m.nextUserID++
	}

	apiKey := m.createAPIKeyFn()
	user.APIKey = apiKey

	m.users[uint(user.ID)] = user
	return apiKey, nil
}

func (m *MockUserRepo) GetByID(ctx context.Context, ID uint) (domain.User, error) {
	if m.getByIDError != nil {
		return domain.User{}, m.getByIDError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	user, exists := m.users[ID]
	if !exists {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepo) UpdateUserBalance(ctx context.Context, ID domain.UserID, amount float64) error {
	if m.updateError != nil {
		return m.updateError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	user, exists := m.users[uint(ID)]
	if !exists {
		return errors.New("user not found")
	}

	user.Balance = amount
	m.users[uint(ID)] = user
	return nil
}

func (m *MockUserRepo) SetUsers(users map[uint]domain.User) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users = users
}

func (m *MockUserRepo) SetCreateError(err error) {
	m.createError = err
}

func (m *MockUserRepo) SetGetByIDError(err error) {
	m.getByIDError = err
}

func (m *MockUserRepo) SetUpdateError(err error) {
	m.updateError = err
}

func (m *MockUserRepo) SetCreateAPIKeyFn(fn func() domain.APIKey) {
	m.createAPIKeyFn = fn
}

type MockRabbit struct {
	PublishedMessages map[string][][]byte
	PublishError      error
	mu                sync.Mutex
	Logger            *slog.Logger
}

func NewMockRabbit() *MockRabbit {
	return &MockRabbit{
		PublishedMessages: make(map[string][][]byte),
		Logger:            slog.Default(),
	}
}

func (m *MockRabbit) Publish(body []byte, queue string) error {
	if m.PublishError != nil {
		return m.PublishError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.PublishedMessages[queue]; !exists {
		m.PublishedMessages[queue] = make([][]byte, 0)
	}

	m.PublishedMessages[queue] = append(m.PublishedMessages[queue], body)

	return nil
}

func (m *MockRabbit) SetPublishError(err error) {
	m.PublishError = err
}

func (m *MockRabbit) ClearPublishedMessages() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.PublishedMessages = make(map[string][][]byte)
}

type MockUserService struct {
	repo   *MockUserRepo
	rabbit *MockRabbit
}

func (s *MockUserService) CreateUser(ctx context.Context, user domain.User) (domain.APIKey, error) {
	return s.repo.Create(ctx, user)
}

func (s *MockUserService) GetUserByID(ctx context.Context, ID domain.UserID) (domain.User, error) {
	return s.repo.GetByID(ctx, uint(ID))
}

func (s *MockUserService) CreditUserBalance(ctx context.Context, ID domain.UserID, amount float64) error {
	user, err := s.repo.GetByID(ctx, uint(ID))
	if err != nil {
		return err
	}
	return s.repo.UpdateUserBalance(ctx, ID, user.Balance+amount)
}

func (s *MockUserService) DebitUserBalance(ctx context.Context, body []byte) (event.SMSUpdateEvent, error) {
	return event.SMSUpdateEvent{
		Domain: event.SMS,
		SMSID:  1,
		Status: event.StatusSuccess,
	}, nil
}

func (s *MockUserService) UnsuccessfulSMS(ctx context.Context, body []byte) error {
	return nil
}

func (s *MockUserService) UpdateSMSStatus(ctx context.Context, sms event.SMSUpdateEvent) error {
	body, err := json.Marshal(sms)
	if err != nil {
		return err
	}
	return s.rabbit.Publish(body, "sms-update-queue")
}
