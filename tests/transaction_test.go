package tests

import (
	"finance/internal/domain/entities"
	"finance/internal/domain/valueobjects"
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTransaction(t *testing.T) {
	t.Run("successful transaction creation for debit", func(t *testing.T) {
		walletID := uuid.New()
		userID := uuid.New()
		smsID := uuid.New()
		amount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		tx := entities.NewTransaction(walletID, userID, smsID, amount, entities.TransactionDebit)

		assert.NotEqual(t, uuid.Nil, tx.ID)
		assert.Equal(t, walletID, tx.WalletID)
		assert.Equal(t, userID, tx.UserID)
		assert.Equal(t, smsID, tx.SMSID)
		assert.Equal(t, amount, tx.Amount)
		assert.Equal(t, entities.TransactionDebit, tx.Type)
		assert.Equal(t, entities.TransactionPending, tx.Status)
		assert.False(t, tx.CreatedAt.IsZero())
		assert.False(t, tx.UpdatedAt.IsZero())
	})

	t.Run("successful transaction creation for credit", func(t *testing.T) {
		walletID := uuid.New()
		userID := uuid.New()
		smsID := uuid.New()
		amount, _ := valueobjects.NewMoney(big.NewInt(200), "EUR")

		tx := entities.NewTransaction(walletID, userID, smsID, amount, entities.TransactionCredit)

		assert.NotEqual(t, uuid.Nil, tx.ID)
		assert.Equal(t, walletID, tx.WalletID)
		assert.Equal(t, userID, tx.UserID)
		assert.Equal(t, smsID, tx.SMSID)
		assert.Equal(t, amount, tx.Amount)
		assert.Equal(t, entities.TransactionCredit, tx.Type)
		assert.Equal(t, entities.TransactionPending, tx.Status)
		assert.False(t, tx.CreatedAt.IsZero())
		assert.False(t, tx.UpdatedAt.IsZero())
	})

	t.Run("each transaction should have unique ID", func(t *testing.T) {
		walletID := uuid.New()
		userID := uuid.New()
		smsID := uuid.New()
		amount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		tx1 := entities.NewTransaction(walletID, userID, smsID, amount, entities.TransactionDebit)
		tx2 := entities.NewTransaction(walletID, userID, smsID, amount, entities.TransactionDebit)

		assert.NotEqual(t, tx1.ID, tx2.ID)
	})
}

func TestTransaction_MarkCompleted(t *testing.T) {
	t.Run("should successfully mark pending transaction as completed", func(t *testing.T) {
		walletID := uuid.New()
		userID := uuid.New()
		smsID := uuid.New()
		amount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		tx := entities.NewTransaction(walletID, userID, smsID, amount, entities.TransactionDebit)
		originalUpdatedAt := tx.UpdatedAt

		err := tx.MarkCompleted()

		require.NoError(t, err)
		assert.Equal(t, entities.TransactionCompleted, tx.Status)
		assert.True(t, tx.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("should fail when transaction is already completed", func(t *testing.T) {
		walletID := uuid.New()
		userID := uuid.New()
		smsID := uuid.New()
		amount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		tx := entities.NewTransaction(walletID, userID, smsID, amount, entities.TransactionDebit)
		tx.MarkCompleted()

		err := tx.MarkCompleted()

		assert.Error(t, err)
		assert.Equal(t, entities.ErrInvalidTransactionState, err)
	})

	t.Run("should fail when transaction is already failed", func(t *testing.T) {
		walletID := uuid.New()
		userID := uuid.New()
		smsID := uuid.New()
		amount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		tx := entities.NewTransaction(walletID, userID, smsID, amount, entities.TransactionDebit)
		tx.MarkFailed()

		err := tx.MarkCompleted()

		assert.Error(t, err)
		assert.Equal(t, entities.ErrInvalidTransactionState, err)
	})
}

func TestTransaction_MarkFailed(t *testing.T) {
	t.Run("should successfully mark pending transaction as failed", func(t *testing.T) {
		walletID := uuid.New()
		userID := uuid.New()
		smsID := uuid.New()
		amount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		tx := entities.NewTransaction(walletID, userID, smsID, amount, entities.TransactionDebit)
		originalUpdatedAt := tx.UpdatedAt

		err := tx.MarkFailed()

		require.NoError(t, err)
		assert.Equal(t, entities.TransactionFailed, tx.Status)
		assert.True(t, tx.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("should fail when transaction is already completed", func(t *testing.T) {
		walletID := uuid.New()
		userID := uuid.New()
		smsID := uuid.New()
		amount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		tx := entities.NewTransaction(walletID, userID, smsID, amount, entities.TransactionDebit)
		tx.MarkCompleted()

		err := tx.MarkFailed()

		assert.Error(t, err)
		assert.Equal(t, entities.ErrInvalidTransactionState, err)
	})

	t.Run("should fail when transaction is already failed", func(t *testing.T) {
		walletID := uuid.New()
		userID := uuid.New()
		smsID := uuid.New()
		amount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		tx := entities.NewTransaction(walletID, userID, smsID, amount, entities.TransactionDebit)
		tx.MarkFailed()

		err := tx.MarkFailed()

		assert.Error(t, err)
		assert.Equal(t, entities.ErrInvalidTransactionState, err)
	})
}

func TestTransactionStatus_Constants(t *testing.T) {
	t.Run("should have correct status constants", func(t *testing.T) {
		assert.Equal(t, entities.TransactionStatus("pending"), entities.TransactionPending)
		assert.Equal(t, entities.TransactionStatus("completed"), entities.TransactionCompleted)
		assert.Equal(t, entities.TransactionStatus("failed"), entities.TransactionFailed)
	})
}

func TestTransactionType_Constants(t *testing.T) {
	t.Run("should have correct type constants", func(t *testing.T) {
		assert.Equal(t, entities.TransactionType("debit"), entities.TransactionDebit)
		assert.Equal(t, entities.TransactionType("credit"), entities.TransactionCredit)
	})
}
