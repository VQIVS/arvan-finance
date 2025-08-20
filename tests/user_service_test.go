package tests

import (
	"billing-service/internal/user/domain"
	"billing-service/internal/user/event"
	"billing-service/internal/user/port"
	"billing-service/pkg/constants"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

type TestableUserService struct {
	service   port.Service
	mockRepo  *MockUserRepo
	mockQueue map[string][][]byte
}

func NewTestableUserService() *TestableUserService {
	mockRepo := NewMockUserRepo()
	service := &testUserService{
		repo:      mockRepo,
		mockQueue: make(map[string][][]byte),
	}

	return &TestableUserService{
		service:   service,
		mockRepo:  mockRepo,
		mockQueue: service.mockQueue,
	}
}

type testUserService struct {
	repo      port.Repo
	mockQueue map[string][][]byte
}

func (s *testUserService) CreateUser(ctx context.Context, user domain.User) (domain.APIKey, error) {
	apiKey, err := s.repo.Create(ctx, user)
	if err != nil {
		return "", err
	}
	return apiKey, nil
}

func (s *testUserService) GetUserByID(ctx context.Context, ID domain.UserID) (domain.User, error) {
	user, err := s.repo.GetByID(ctx, uint(ID))
	if err != nil {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

func (s *testUserService) CreditUserBalance(ctx context.Context, ID domain.UserID, amount float64) error {
	if amount <= 0 {
		return errors.New("invalid credit amount")
	}

	user, err := s.repo.GetByID(ctx, uint(ID))
	if err != nil {
		return errors.New("user not found")
	}

	user.Balance += amount
	return s.repo.UpdateUserBalance(ctx, ID, user.Balance)
}

func (s *testUserService) DebitUserBalance(ctx context.Context, body []byte) (event.SMSUpdateEvent, error) {
	var msg event.UserBalanceEvent
	if err := json.Unmarshal(body, &msg); err != nil {
		return event.SMSUpdateEvent{}, err
	}
	if msg.Amount <= 0 {
		return event.SMSUpdateEvent{}, errors.New("invalid debit amount")
	}

	user, err := s.repo.GetByID(ctx, msg.UserID)
	if err != nil {
		return event.SMSUpdateEvent{}, errors.New("user not found")
	}
	if user.Balance < msg.Amount {
		return event.SMSUpdateEvent{}, errors.New("insufficient balance")
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

func (s *testUserService) UnsuccessfulSMS(ctx context.Context, body []byte) error {
	var msg event.UserBalanceEvent
	if err := json.Unmarshal(body, &msg); err != nil {
		return err
	}

	user, err := s.repo.GetByID(ctx, msg.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	if msg.Amount <= 0 {
		return errors.New("invalid SMS amount")
	}

	return s.repo.UpdateUserBalance(ctx, domain.UserID(msg.UserID), user.Balance+msg.Amount)
}

func (s *testUserService) UpdateSMSStatus(ctx context.Context, sms event.SMSUpdateEvent) error {
	body, err := json.Marshal(sms)
	if err != nil {
		return err
	}

	if _, exists := s.mockQueue[constants.KeySMSUpdate]; !exists {
		s.mockQueue[constants.KeySMSUpdate] = make([][]byte, 0)
	}
	s.mockQueue[constants.KeySMSUpdate] = append(s.mockQueue[constants.KeySMSUpdate], body)

	return nil
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		user           domain.User
		repoError      error
		expectedError  error
		expectedAPIKey domain.APIKey
	}{
		{
			name: "Success",
			user: domain.User{
				Balance:   100.0,
				CreatedAt: time.Now(),
			},
			expectedAPIKey: "test-api-key",
			expectedError:  nil,
		},
		{
			name: "Repository Error",
			user: domain.User{
				Balance:   100.0,
				CreatedAt: time.Now(),
			},
			repoError:     errors.New("repository error"),
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTestableUserService()

			if tt.repoError != nil {
				service.mockRepo.SetCreateError(tt.repoError)
			}

			apiKey, err := service.service.CreateUser(context.Background(), tt.user)

			if (err != nil) != (tt.expectedError != nil) {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			if err == nil && apiKey != tt.expectedAPIKey {
				t.Errorf("Expected API key: %v, got: %v", tt.expectedAPIKey, apiKey)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        domain.UserID
		mockUser      domain.User
		setupMock     func(*MockUserRepo)
		expectedError bool
	}{
		{
			name:   "Success",
			userID: 1,
			mockUser: domain.User{
				ID:        1,
				Balance:   100.0,
				APIKey:    "test-api-key",
				CreatedAt: time.Now(),
			},
			setupMock: func(repo *MockUserRepo) {
				users := make(map[uint]domain.User)
				users[1] = domain.User{
					ID:        1,
					Balance:   100.0,
					APIKey:    "test-api-key",
					CreatedAt: time.Now(),
				}
				repo.SetUsers(users)
			},
			expectedError: false,
		},
		{
			name:   "User Not Found",
			userID: 999,
			setupMock: func(repo *MockUserRepo) {
			},
			expectedError: true,
		},
		{
			name:   "Repository Error",
			userID: 1,
			setupMock: func(repo *MockUserRepo) {
				repo.SetGetByIDError(errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTestableUserService()

			if tt.setupMock != nil {
				tt.setupMock(service.mockRepo)
			}

			got, err := service.service.GetUserByID(ctx, tt.userID)

			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			if err == nil && got.ID != tt.mockUser.ID {
				t.Errorf("Expected user ID: %v, got: %v", tt.mockUser.ID, got.ID)
			}

			if err == nil && got.APIKey != tt.mockUser.APIKey {
				t.Errorf("Expected user APIKey: %v, got: %v", tt.mockUser.APIKey, got.APIKey)
			}

			if err == nil && got.Balance != tt.mockUser.Balance {
				t.Errorf("Expected user Balance: %v, got: %v", tt.mockUser.Balance, got.Balance)
			}
		})
	}
}

func TestCreditUserBalance(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name            string
		userID          domain.UserID
		amount          float64
		setupMock       func(*MockUserRepo)
		expectedError   bool
		expectedBalance float64
	}{
		{
			name:   "Success",
			userID: 1,
			amount: 50.0,
			setupMock: func(repo *MockUserRepo) {
				users := make(map[uint]domain.User)
				users[1] = domain.User{
					ID:        1,
					Balance:   100.0,
					APIKey:    "test-api-key",
					CreatedAt: time.Now(),
				}
				repo.SetUsers(users)
			},
			expectedError:   false,
			expectedBalance: 150.0,
		},
		{
			name:   "Negative Amount",
			userID: 1,
			amount: -50.0,
			setupMock: func(repo *MockUserRepo) {
				users := make(map[uint]domain.User)
				users[1] = domain.User{
					ID:        1,
					Balance:   100.0,
					APIKey:    "test-api-key",
					CreatedAt: time.Now(),
				}
				repo.SetUsers(users)
			},
			expectedError: true,
		},
		{
			name:   "User Not Found",
			userID: 999,
			amount: 50.0,
			setupMock: func(repo *MockUserRepo) {
			},
			expectedError: true,
		},
		{
			name:   "Update Error",
			userID: 1,
			amount: 50.0,
			setupMock: func(repo *MockUserRepo) {
				users := make(map[uint]domain.User)
				users[1] = domain.User{
					ID:        1,
					Balance:   100.0,
					APIKey:    "test-api-key",
					CreatedAt: time.Now(),
				}
				repo.SetUsers(users)
				repo.SetUpdateError(errors.New("update error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTestableUserService()

			if tt.setupMock != nil {
				tt.setupMock(service.mockRepo)
			}

			err := service.service.CreditUserBalance(ctx, tt.userID, tt.amount)

			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			if err == nil {
				updatedUser, _ := service.mockRepo.GetByID(ctx, uint(tt.userID))
				if updatedUser.Balance != tt.expectedBalance {
					t.Errorf("Expected balance: %v, got: %v", tt.expectedBalance, updatedUser.Balance)
				}
			}
		})
	}
}

func TestDebitUserBalance(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		event          event.UserBalanceEvent
		setupMock      func(*MockUserRepo)
		expectedError  bool
		expectedStatus event.Status
	}{
		{
			name: "Success",
			event: event.UserBalanceEvent{
				UserID: 1,
				SMSID:  123,
				Amount: 50.0,
				Type:   event.SMSDebitEvent,
				Domain: event.SMS,
			},
			setupMock: func(repo *MockUserRepo) {
				users := make(map[uint]domain.User)
				users[1] = domain.User{
					ID:        1,
					Balance:   100.0,
					APIKey:    "test-api-key",
					CreatedAt: time.Now(),
				}
				repo.SetUsers(users)
			},
			expectedError:  false,
			expectedStatus: event.StatusSuccess,
		},
		{
			name: "Invalid Amount",
			event: event.UserBalanceEvent{
				UserID: 1,
				SMSID:  123,
				Amount: -50.0,
				Type:   event.SMSDebitEvent,
				Domain: event.SMS,
			},
			setupMock: func(repo *MockUserRepo) {
				users := make(map[uint]domain.User)
				users[1] = domain.User{
					ID:        1,
					Balance:   100.0,
					APIKey:    "test-api-key",
					CreatedAt: time.Now(),
				}
				repo.SetUsers(users)
			},
			expectedError: true,
		},
		{
			name: "User Not Found",
			event: event.UserBalanceEvent{
				UserID: 999,
				SMSID:  123,
				Amount: 50.0,
				Type:   event.SMSDebitEvent,
				Domain: event.SMS,
			},
			setupMock: func(repo *MockUserRepo) {
			},
			expectedError: true,
		},
		{
			name: "Insufficient Balance",
			event: event.UserBalanceEvent{
				UserID: 1,
				SMSID:  123,
				Amount: 200.0,
				Type:   event.SMSDebitEvent,
				Domain: event.SMS,
			},
			setupMock: func(repo *MockUserRepo) {
				users := make(map[uint]domain.User)
				users[1] = domain.User{
					ID:        1,
					Balance:   100.0,
					APIKey:    "test-api-key",
					CreatedAt: time.Now(),
				}
				repo.SetUsers(users)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTestableUserService()

			if tt.setupMock != nil {
				tt.setupMock(service.mockRepo)
			}

			eventBytes, _ := json.Marshal(tt.event)

			result, err := service.service.DebitUserBalance(ctx, eventBytes)

			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
				return
			}

			if err == nil {
				if result.Status != tt.expectedStatus {
					t.Errorf("Expected status: %v, got: %v", tt.expectedStatus, result.Status)
				}

				if result.SMSID != tt.event.SMSID {
					t.Errorf("Expected SMS ID: %v, got: %v", tt.event.SMSID, result.SMSID)
				}

				updatedUser, _ := service.mockRepo.GetByID(ctx, tt.event.UserID)
				expectedBalance := 100.0 - tt.event.Amount
				if updatedUser.Balance != expectedBalance {
					t.Errorf("Expected balance: %v, got: %v", expectedBalance, updatedUser.Balance)
				}
			}
		})
	}
}

func TestUnsuccessfulSMS(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name            string
		event           event.UserBalanceEvent
		setupMock       func(*MockUserRepo)
		expectedError   bool
		expectedBalance float64
	}{
		{
			name: "Success",
			event: event.UserBalanceEvent{
				UserID: 1,
				SMSID:  123,
				Amount: 50.0,
				Type:   event.SMSCreditEvent,
				Domain: event.SMS,
			},
			setupMock: func(repo *MockUserRepo) {
				users := make(map[uint]domain.User)
				users[1] = domain.User{
					ID:        1,
					Balance:   100.0,
					APIKey:    "test-api-key",
					CreatedAt: time.Now(),
				}
				repo.SetUsers(users)
			},
			expectedError:   false,
			expectedBalance: 150.0, // 100 + 50
		},
		{
			name: "Invalid Amount",
			event: event.UserBalanceEvent{
				UserID: 1,
				SMSID:  123,
				Amount: -50.0,
				Type:   event.SMSCreditEvent,
				Domain: event.SMS,
			},
			setupMock: func(repo *MockUserRepo) {
				users := make(map[uint]domain.User)
				users[1] = domain.User{
					ID:        1,
					Balance:   100.0,
					APIKey:    "test-api-key",
					CreatedAt: time.Now(),
				}
				repo.SetUsers(users)
			},
			expectedError: true,
		},
		{
			name: "User Not Found",
			event: event.UserBalanceEvent{
				UserID: 999,
				SMSID:  123,
				Amount: 50.0,
				Type:   event.SMSCreditEvent,
				Domain: event.SMS,
			},
			setupMock: func(repo *MockUserRepo) {
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTestableUserService()

			if tt.setupMock != nil {
				tt.setupMock(service.mockRepo)
			}

			eventBytes, _ := json.Marshal(tt.event)

			err := service.service.UnsuccessfulSMS(ctx, eventBytes)

			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
				return
			}

			if err == nil {
				updatedUser, _ := service.mockRepo.GetByID(ctx, tt.event.UserID)
				if updatedUser.Balance != tt.expectedBalance {
					t.Errorf("Expected balance: %v, got: %v", tt.expectedBalance, updatedUser.Balance)
				}
			}
		})
	}
}

func TestUpdateSMSStatus(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		event         event.SMSUpdateEvent
		expectedError bool
		expectedQueue string
	}{
		{
			name: "Success",
			event: event.SMSUpdateEvent{
				SMSID:  123,
				Status: event.StatusSuccess,
				Domain: event.SMS,
			},
			expectedError: false,
			expectedQueue: constants.KeySMSUpdate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTestableUserService()

			err := service.service.UpdateSMSStatus(ctx, tt.event)

			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
				return
			}

			if err == nil {
				if len(service.mockQueue[tt.expectedQueue]) == 0 {
					t.Errorf("Expected message to be published to queue %s, but none found", tt.expectedQueue)
				}

				var publishedEvent event.SMSUpdateEvent
				json.Unmarshal(service.mockQueue[tt.expectedQueue][0], &publishedEvent)

				if publishedEvent.SMSID != tt.event.SMSID {
					t.Errorf("Expected SMS ID: %v, got: %v", tt.event.SMSID, publishedEvent.SMSID)
				}

				if publishedEvent.Status != tt.event.Status {
					t.Errorf("Expected Status: %v, got: %v", tt.event.Status, publishedEvent.Status)
				}
			}
		})
	}
}
