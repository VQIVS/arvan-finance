package tests

import (
	"context"
	"errors"
	"finance/internal/domain/entities"
	"finance/internal/domain/events"
	"finance/internal/domain/valueobjects"
	"finance/internal/usecase"
	"finance/pkg/logger"
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type MockWalletRepo struct {
	mock.Mock
}

func (m *MockWalletRepo) Save(ctx context.Context, wallet *entities.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepo) FindByID(ctx context.Context, ID uuid.UUID) (*entities.Wallet, error) {
	args := m.Called(ctx, ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Wallet), args.Error(1)
}

func (m *MockWalletRepo) FindByUserID(ctx context.Context, userID uuid.UUID) (*entities.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Wallet), args.Error(1)
}

func (m *MockWalletRepo) UpdateBalance(ctx context.Context, wallet *entities.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepo) WithTx(tx *gorm.DB) entities.WalletRepo {
	args := m.Called(tx)
	return args.Get(0).(entities.WalletRepo)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByID(ctx context.Context, ID uuid.UUID) (*entities.User, error) {
	args := m.Called(ctx, ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepo) GetAll(ctx context.Context) ([]*entities.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserRepo) WithTx(tx *gorm.DB) entities.UserRepo {
	args := m.Called(tx)
	return args.Get(0).(entities.UserRepo)
}

type MockTransactionRepo struct {
	mock.Mock
}

func (m *MockTransactionRepo) Create(ctx context.Context, tx *entities.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTransactionRepo) FindByID(ctx context.Context, id string) (*entities.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Transaction), args.Error(1)
}

func (m *MockTransactionRepo) UpdateStatus(ctx context.Context, tx *entities.Transaction, status entities.TransactionStatus) error {
	args := m.Called(ctx, tx, status)
	return args.Error(0)
}

func (m *MockTransactionRepo) WithTx(tx *gorm.DB) entities.TransactionRepo {
	args := m.Called(tx)
	return args.Get(0).(entities.TransactionRepo)
}

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) WithTransaction(fn func(tx *gorm.DB) error) error {
	args := m.Called(fn)
	if fn != nil {
		return fn(nil)
	}
	return args.Error(0)
}

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) PublishEvent(ctx context.Context, event events.SMSEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func setupWalletServiceTest() (*usecase.WalletService, *MockWalletRepo, *MockUserRepo, *MockTransactionRepo, *MockTransactionManager, *MockPublisher) {
	mockWalletRepo := &MockWalletRepo{}
	mockUserRepo := &MockUserRepo{}
	mockTransactionRepo := &MockTransactionRepo{}
	mockTxManager := &MockTransactionManager{}
	mockPublisher := &MockPublisher{}
	mockLogger := &logger.Logger{}

	mockWalletRepo.On("WithTx", mock.Anything).Return(mockWalletRepo)
	mockUserRepo.On("WithTx", mock.Anything).Return(mockUserRepo)
	mockTransactionRepo.On("WithTx", mock.Anything).Return(mockTransactionRepo)

	service := usecase.NewWalletService(
		mockWalletRepo,
		mockUserRepo,
		mockTransactionRepo,
		mockTxManager,
		mockPublisher,
		mockLogger,
	)

	return service, mockWalletRepo, mockUserRepo, mockTransactionRepo, mockTxManager, mockPublisher
}

func TestWalletService_DebitUserBalance(t *testing.T) {
	t.Run("successful debit operation", func(t *testing.T) {
		service, mockWalletRepo, _, mockTransactionRepo, mockTxManager, _ := setupWalletServiceTest()

		userID := uuid.New()
		smsID := uuid.New()
		amount := *big.NewInt(100)
		ctx := context.Background()

		wallet, _ := entities.NewWallet(userID, "USD")
		initialAmount, _ := valueobjects.NewMoney(big.NewInt(200), "USD")
		wallet.Credit(initialAmount)

		mockTxManager.On("WithTransaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(nil)
		mockWalletRepo.On("FindByUserID", ctx, userID).Return(wallet, nil)
		mockTransactionRepo.On("Create", ctx, mock.AnythingOfType("*entities.Transaction")).Return(nil)
		mockWalletRepo.On("UpdateBalance", ctx, mock.AnythingOfType("*entities.Wallet")).Return(nil)
		mockTransactionRepo.On("UpdateStatus", ctx, mock.AnythingOfType("*entities.Transaction"), entities.TransactionCompleted).Return(nil)

		event, err := service.DebitUserbalance(ctx, userID, smsID, amount)

		require.NoError(t, err)
		require.NotNil(t, event)
		assert.Equal(t, userID.String(), event.UserID)
		assert.Equal(t, smsID.String(), event.SMSID)
		assert.Equal(t, int64(100), event.Amount)
		assert.NotEmpty(t, event.TransactionID)

		mockWalletRepo.AssertExpectations(t)
		mockTransactionRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("should fail when wallet not found", func(t *testing.T) {
		service, mockWalletRepo, _, _, mockTxManager, _ := setupWalletServiceTest()

		userID := uuid.New()
		smsID := uuid.New()
		amount := *big.NewInt(100)
		ctx := context.Background()

		mockTxManager.On("WithTransaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(entities.ErrWalletNotFound)
		mockWalletRepo.On("FindByUserID", ctx, userID).Return(nil, entities.ErrWalletNotFound)

		event, err := service.DebitUserbalance(ctx, userID, smsID, amount)

		assert.Error(t, err)
		assert.Nil(t, event)
		assert.Equal(t, entities.ErrWalletNotFound, err)

		mockWalletRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("should fail when insufficient balance", func(t *testing.T) {
		service, mockWalletRepo, _, mockTransactionRepo, mockTxManager, _ := setupWalletServiceTest()

		userID := uuid.New()
		smsID := uuid.New()
		amount := *big.NewInt(100)
		ctx := context.Background()

		wallet, _ := entities.NewWallet(userID, "USD")

		mockTxManager.On("WithTransaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(entities.ErrInsufficientBalance)
		mockWalletRepo.On("FindByUserID", ctx, userID).Return(wallet, nil)
		mockTransactionRepo.On("Create", ctx, mock.AnythingOfType("*entities.Transaction")).Return(nil)

		event, err := service.DebitUserbalance(ctx, userID, smsID, amount)

		assert.Error(t, err)
		assert.Nil(t, event)
		assert.Equal(t, entities.ErrInsufficientBalance, err)

		mockWalletRepo.AssertExpectations(t)
		mockTransactionRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})
}

func TestWalletService_CreditUserBalance(t *testing.T) {
	t.Run("successful credit operation", func(t *testing.T) {
		service, mockWalletRepo, _, mockTransactionRepo, mockTxManager, _ := setupWalletServiceTest()

		userID := uuid.New()
		amount := *big.NewInt(100)
		ctx := context.Background()

		wallet, _ := entities.NewWallet(userID, "USD")

		mockTxManager.On("WithTransaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(nil)
		mockWalletRepo.On("FindByUserID", ctx, userID).Return(wallet, nil)
		mockTransactionRepo.On("Create", ctx, mock.AnythingOfType("*entities.Transaction")).Return(nil)
		mockWalletRepo.On("UpdateBalance", ctx, mock.AnythingOfType("*entities.Wallet")).Return(nil)
		mockTransactionRepo.On("UpdateStatus", ctx, mock.AnythingOfType("*entities.Transaction"), entities.TransactionCompleted).Return(nil)

		err := service.CreditUserBalance(ctx, userID, amount)

		require.NoError(t, err)

		mockWalletRepo.AssertExpectations(t)
		mockTransactionRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("should fail when wallet not found", func(t *testing.T) {
		service, mockWalletRepo, _, _, mockTxManager, _ := setupWalletServiceTest()

		userID := uuid.New()
		amount := *big.NewInt(100)
		ctx := context.Background()

		mockTxManager.On("WithTransaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(entities.ErrWalletNotFound)
		mockWalletRepo.On("FindByUserID", ctx, userID).Return(nil, entities.ErrWalletNotFound)

		err := service.CreditUserBalance(ctx, userID, amount)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrWalletNotFound, err)

		mockWalletRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})
}

func TestWalletService_GetWalletByUserID(t *testing.T) {
	t.Run("successful wallet retrieval", func(t *testing.T) {
		mockWalletRepo := &MockWalletRepo{}
		mockLogger := &logger.Logger{}

		service := usecase.NewWalletService(
			mockWalletRepo,
			nil,
			nil,
			nil,
			nil,
			mockLogger,
		)

		userID := uuid.New()
		ctx := context.Background()
		expectedWallet, _ := entities.NewWallet(userID, "USD")

		mockWalletRepo.On("FindByUserID", ctx, userID).Return(expectedWallet, nil)

		wallet, err := service.GetWalletByUserID(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, expectedWallet, wallet)

		mockWalletRepo.AssertExpectations(t)
	})

	t.Run("should return error when wallet not found", func(t *testing.T) {
		mockWalletRepo := &MockWalletRepo{}
		mockLogger := &logger.Logger{}

		service := usecase.NewWalletService(
			mockWalletRepo,
			nil,
			nil,
			nil,
			nil,
			mockLogger,
		)

		userID := uuid.New()
		ctx := context.Background()

		mockWalletRepo.On("FindByUserID", ctx, userID).Return(nil, entities.ErrWalletNotFound)

		wallet, err := service.GetWalletByUserID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, wallet)
		assert.Equal(t, entities.ErrWalletNotFound, err)

		mockWalletRepo.AssertExpectations(t)
	})
}

func TestWalletService_GetUserByID(t *testing.T) {
	t.Run("successful user retrieval", func(t *testing.T) {
		mockUserRepo := &MockUserRepo{}
		mockLogger := &logger.Logger{}

		service := usecase.NewWalletService(
			nil,
			mockUserRepo,
			nil,
			nil,
			nil,
			mockLogger,
		)

		userID := uuid.New()
		ctx := context.Background()
		expectedUser := &entities.User{
			ID:       userID,
			Name:     "John",
			LastName: "Doe",
			Phone:    "+1234567890",
		}

		mockUserRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)

		user, err := service.GetUserByID(ctx, userID)

		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestWalletService_GetAllUsers(t *testing.T) {
	t.Run("successful users retrieval", func(t *testing.T) {
		mockUserRepo := &MockUserRepo{}
		mockLogger := &logger.Logger{}

		service := usecase.NewWalletService(
			nil,
			mockUserRepo,
			nil,
			nil,
			nil,
			mockLogger,
		)

		ctx := context.Background()
		expectedUsers := []*entities.User{
			{
				ID:       uuid.New(),
				Name:     "John",
				LastName: "Doe",
				Phone:    "+1234567890",
			},
			{
				ID:       uuid.New(),
				Name:     "Jane",
				LastName: "Smith",
				Phone:    "+0987654321",
			},
		}

		mockUserRepo.On("GetAll", ctx).Return(expectedUsers, nil)

		users, err := service.GetAllUsers(ctx)

		require.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Len(t, users, 2)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestWalletService_RefundTransaction(t *testing.T) {
	t.Run("successful refund operation", func(t *testing.T) {
		service, mockWalletRepo, _, mockTransactionRepo, mockTxManager, _ := setupWalletServiceTest()

		txID := uuid.New().String()
		userID := uuid.New()
		walletID := uuid.New()
		smsID := uuid.New()
		ctx := context.Background()

		amount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		originalTx := entities.NewTransaction(walletID, userID, smsID, amount, entities.TransactionDebit)
		originalTx.MarkCompleted()

		wallet, _ := entities.NewWallet(userID, "USD")

		mockTxManager.On("WithTransaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(nil)
		mockTransactionRepo.On("FindByID", ctx, txID).Return(originalTx, nil)
		mockWalletRepo.On("FindByID", ctx, walletID).Return(wallet, nil)
		mockTransactionRepo.On("Create", ctx, mock.AnythingOfType("*entities.Transaction")).Return(nil)
		mockWalletRepo.On("UpdateBalance", ctx, mock.AnythingOfType("*entities.Wallet")).Return(nil)
		mockTransactionRepo.On("UpdateStatus", ctx, mock.AnythingOfType("*entities.Transaction"), entities.TransactionCompleted).Return(nil)

		err := service.RefundTransaction(ctx, txID)

		require.NoError(t, err)

		mockWalletRepo.AssertExpectations(t)
		mockTransactionRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("should succeed but not process credit transactions", func(t *testing.T) {
		service, _, _, mockTransactionRepo, mockTxManager, _ := setupWalletServiceTest()

		txID := uuid.New().String()
		userID := uuid.New()
		walletID := uuid.New()
		smsID := uuid.New()
		ctx := context.Background()

		amount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		originalTx := entities.NewTransaction(walletID, userID, smsID, amount, entities.TransactionCredit)

		mockTxManager.On("WithTransaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(nil)
		mockTransactionRepo.On("FindByID", ctx, txID).Return(originalTx, nil)

		err := service.RefundTransaction(ctx, txID)

		require.NoError(t, err)

		mockTransactionRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})

	t.Run("should fail when transaction not found", func(t *testing.T) {
		service, _, _, mockTransactionRepo, mockTxManager, _ := setupWalletServiceTest()

		txID := uuid.New().String()
		ctx := context.Background()

		mockTxManager.On("WithTransaction", mock.AnythingOfType("func(*gorm.DB) error")).Return(errors.New("transaction not found"))
		mockTransactionRepo.On("FindByID", ctx, txID).Return(nil, errors.New("transaction not found"))

		err := service.RefundTransaction(ctx, txID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction not found")

		mockTransactionRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})
}
