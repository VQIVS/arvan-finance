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

func TestIntegration_WalletOperations(t *testing.T) {
	t.Run("complete wallet lifecycle", func(t *testing.T) {
		userID := uuid.New()
		wallet, err := entities.NewWallet(userID, "USD")
		require.NoError(t, err)

		assert.Equal(t, userID, wallet.UserID)
		assert.Equal(t, "USD", wallet.Currency)
		assert.True(t, wallet.Balance.IsZero())

		creditAmount, err := valueobjects.NewMoney(big.NewInt(1000), "USD")
		require.NoError(t, err)

		err = wallet.Credit(creditAmount)
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(1000), wallet.Balance.Amount())

		smsID := uuid.New()
		debitAmount, err := valueobjects.NewMoney(big.NewInt(250), "USD")
		require.NoError(t, err)

		err = wallet.HasSufficientBalance(debitAmount)
		require.NoError(t, err)

		transaction := entities.NewTransaction(wallet.ID, userID, smsID, debitAmount, entities.TransactionDebit)
		assert.Equal(t, entities.TransactionPending, transaction.Status)
		assert.Equal(t, entities.TransactionDebit, transaction.Type)

		err = wallet.Debit(debitAmount)
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(750), wallet.Balance.Amount())

		err = transaction.MarkCompleted()
		require.NoError(t, err)
		assert.Equal(t, entities.TransactionCompleted, transaction.Status)

		refundTransaction := entities.NewTransaction(wallet.ID, userID, smsID, debitAmount, entities.TransactionCredit)

		err = wallet.Credit(debitAmount)
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(1000), wallet.Balance.Amount())

		err = refundTransaction.MarkCompleted()
		require.NoError(t, err)
		assert.Equal(t, entities.TransactionCompleted, refundTransaction.Status)
	})

	t.Run("error scenarios", func(t *testing.T) {
		userID := uuid.New()
		wallet, _ := entities.NewWallet(userID, "USD")

		debitAmount, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		err := wallet.Debit(debitAmount)
		assert.Error(t, err)
		assert.Equal(t, entities.ErrInsufficientBalance, err)

		wrongCurrencyAmount, _ := valueobjects.NewMoney(big.NewInt(100), "EUR")
		err = wallet.Credit(wrongCurrencyAmount)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "currency mismatch")

		err = wallet.HasSufficientBalance(wrongCurrencyAmount)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "currency mismatch")

		smsID := uuid.New()
		transaction := entities.NewTransaction(wallet.ID, userID, smsID, debitAmount, entities.TransactionDebit)
		err = transaction.MarkCompleted()
		require.NoError(t, err)

		err = transaction.MarkCompleted()
		assert.Error(t, err)
		assert.Equal(t, entities.ErrInvalidTransactionState, err)
	})
}

func TestIntegration_MoneyOperations(t *testing.T) {
	t.Run("money arithmetic operations", func(t *testing.T) {
		amount1, err := valueobjects.NewMoney(big.NewInt(100), "USD")
		require.NoError(t, err)

		amount2, err := valueobjects.NewMoney(big.NewInt(50), "USD")
		require.NoError(t, err)

		sum, err := amount1.Add(amount2)
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(150), sum.Amount())
		assert.Equal(t, "USD", sum.Currency())

		difference, err := amount1.Subtract(amount2)
		require.NoError(t, err)
		assert.Equal(t, big.NewInt(50), difference.Amount())

		isGreater, err := amount1.GreaterThanOrEqual(amount2)
		require.NoError(t, err)
		assert.True(t, isGreater)

		isEqual, err := amount1.GreaterThanOrEqual(amount1)
		require.NoError(t, err)
		assert.True(t, isEqual)

		isLess, err := amount2.GreaterThanOrEqual(amount1)
		require.NoError(t, err)
		assert.False(t, isLess)
	})

	t.Run("edge cases and validations", func(t *testing.T) {
		zeroAmount, err := valueobjects.NewMoney(big.NewInt(0), "USD")
		require.NoError(t, err)
		assert.True(t, zeroAmount.IsZero())
		assert.False(t, zeroAmount.IsNegative())

		amount, err := valueobjects.NewMoney(big.NewInt(100), "usd")
		require.NoError(t, err)
		assert.Equal(t, "USD", amount.Currency())

		_, err = valueobjects.NewMoney(big.NewInt(-100), "USD")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "amount cannot be negative")

		_, err = valueobjects.NewMoney(nil, "USD")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "amount cannot be nil")
	})
}
