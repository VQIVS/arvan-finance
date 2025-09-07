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

func TestNewWallet(t *testing.T) {
	t.Run("successful wallet creation", func(t *testing.T) {
		userID := uuid.New()
		currency := "USD"

		wallet, err := entities.NewWallet(userID, currency)

		require.NoError(t, err)
		assert.Equal(t, userID, wallet.UserID)
		assert.Equal(t, currency, wallet.Currency)
		assert.True(t, wallet.Balance.IsZero())
		assert.NotEqual(t, uuid.Nil, wallet.ID)
		assert.False(t, wallet.CreatedAt.IsZero())
		assert.False(t, wallet.UpdatedAt.IsZero())
	})

	t.Run("should fail with invalid currency for money creation", func(t *testing.T) {
		userID := uuid.New()
		wallet, err := entities.NewWallet(userID, "")

		require.NoError(t, err)
		assert.Equal(t, "", wallet.Currency)
	})
}

func TestWallet_HasSufficientBalance(t *testing.T) {
	t.Run("should return nil when wallet has sufficient balance", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")

		creditAmount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		wallet.Credit(creditAmount)

		checkAmount, _ := valueobjects.NewMoney(big.NewInt(50), "USD")
		err := wallet.HasSufficientBalance(checkAmount)

		assert.NoError(t, err)
	})

	t.Run("should return error when wallet has insufficient balance", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")

		checkAmount, _ := valueobjects.NewMoney(big.NewInt(50), "USD")
		err := wallet.HasSufficientBalance(checkAmount)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrInsufficientBalance, err)
	})

	t.Run("should return error for currency mismatch", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")

		checkAmount, _ := valueobjects.NewMoney(big.NewInt(50), "EUR")
		err := wallet.HasSufficientBalance(checkAmount)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "currency mismatch")
	})

	t.Run("should return nil when amounts are exactly equal", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")

		creditAmount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		wallet.Credit(creditAmount)

		checkAmount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		err := wallet.HasSufficientBalance(checkAmount)

		assert.NoError(t, err)
	})
}

func TestWallet_Credit(t *testing.T) {
	t.Run("successful credit operation", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")
		originalUpdatedAt := wallet.UpdatedAt

		creditAmount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		err := wallet.Credit(creditAmount)

		require.NoError(t, err)
		assert.Equal(t, big.NewInt(100), wallet.Balance.Amount())
		assert.True(t, wallet.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("should fail with nil wallet", func(t *testing.T) {
		var wallet *entities.Wallet
		creditAmount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		err := wallet.Credit(creditAmount)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrWalletNotFound, err)
	})

	t.Run("should fail with zero amount", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")

		creditAmount, _ := valueobjects.NewMoney(big.NewInt(0), "USD")
		err := wallet.Credit(creditAmount)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrInvalidAmount, err)
	})

	t.Run("should fail with currency mismatch", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")

		creditAmount, _ := valueobjects.NewMoney(big.NewInt(100), "EUR")
		err := wallet.Credit(creditAmount)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "currency mismatch")
	})

	t.Run("multiple credits should accumulate", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")

		creditAmount1, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		creditAmount2, _ := valueobjects.NewMoney(big.NewInt(50), "USD")

		err1 := wallet.Credit(creditAmount1)
		err2 := wallet.Credit(creditAmount2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.Equal(t, big.NewInt(150), wallet.Balance.Amount())
	})
}

func TestWallet_Debit(t *testing.T) {
	t.Run("successful debit operation", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")

		creditAmount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		wallet.Credit(creditAmount)
		originalUpdatedAt := wallet.UpdatedAt

		debitAmount, _ := valueobjects.NewMoney(big.NewInt(30), "USD")
		err := wallet.Debit(debitAmount)

		require.NoError(t, err)
		assert.Equal(t, big.NewInt(70), wallet.Balance.Amount())
		assert.True(t, wallet.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("should fail with nil wallet", func(t *testing.T) {
		var wallet *entities.Wallet
		debitAmount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		err := wallet.Debit(debitAmount)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrWalletNotFound, err)
	})

	t.Run("should fail with zero amount", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")

		debitAmount, _ := valueobjects.NewMoney(big.NewInt(0), "USD")
		err := wallet.Debit(debitAmount)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrInvalidAmount, err)
	})

	t.Run("should fail with insufficient balance", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")

		debitAmount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		err := wallet.Debit(debitAmount)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrInsufficientBalance, err)
	})

	t.Run("should fail when debit would result in negative balance", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")

		creditAmount, _ := valueobjects.NewMoney(big.NewInt(50), "USD")
		wallet.Credit(creditAmount)

		debitAmount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		err := wallet.Debit(debitAmount)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrInsufficientBalance, err)
	})

	t.Run("should allow debit of exact balance", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")

		creditAmount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		wallet.Credit(creditAmount)

		debitAmount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		err := wallet.Debit(debitAmount)

		require.NoError(t, err)
		assert.True(t, wallet.Balance.IsZero())
	})
}
