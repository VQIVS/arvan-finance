package tests

import (
	"finance/internal/domain/valueobjects"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMoney(t *testing.T) {
	t.Run("successful creation with valid amount", func(t *testing.T) {
		amount := big.NewInt(100)
		money, err := valueobjects.NewMoney(amount, "USD")

		require.NoError(t, err)
		assert.Equal(t, "USD", money.Currency())
		assert.Equal(t, big.NewInt(100), money.Amount())
	})

	t.Run("should uppercase currency", func(t *testing.T) {
		amount := big.NewInt(100)
		money, err := valueobjects.NewMoney(amount, "usd")

		require.NoError(t, err)
		assert.Equal(t, "USD", money.Currency())
	})

	t.Run("should fail with nil amount", func(t *testing.T) {
		_, err := valueobjects.NewMoney(nil, "USD")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "amount cannot be nil")
	})

	t.Run("should fail with negative amount", func(t *testing.T) {
		amount := big.NewInt(-100)
		_, err := valueobjects.NewMoney(amount, "USD")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "amount cannot be negative")
	})

	t.Run("should accept zero amount", func(t *testing.T) {
		amount := big.NewInt(0)
		money, err := valueobjects.NewMoney(amount, "USD")

		require.NoError(t, err)
		assert.True(t, money.IsZero())
	})
}

func TestMoney_Add(t *testing.T) {
	t.Run("successful addition with same currency", func(t *testing.T) {
		money1, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		money2, _ := valueobjects.NewMoney(big.NewInt(50), "USD")

		result, err := money1.Add(money2)

		require.NoError(t, err)
		assert.Equal(t, big.NewInt(150), result.Amount())
		assert.Equal(t, "USD", result.Currency())
	})

	t.Run("should fail with different currencies", func(t *testing.T) {
		money1, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		money2, _ := valueobjects.NewMoney(big.NewInt(50), "EUR")

		_, err := money1.Add(money2)

		assert.Error(t, err)
		assert.Equal(t, valueobjects.ErrCurrencyMismatch, err)
	})

	t.Run("adding zero should work", func(t *testing.T) {
		money1, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		money2, _ := valueobjects.NewMoney(big.NewInt(0), "USD")

		result, err := money1.Add(money2)

		require.NoError(t, err)
		assert.Equal(t, big.NewInt(100), result.Amount())
	})
}

func TestMoney_Subtract(t *testing.T) {
	t.Run("successful subtraction with same currency", func(t *testing.T) {
		money1, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		money2, _ := valueobjects.NewMoney(big.NewInt(30), "USD")

		result, err := money1.Subtract(money2)

		require.NoError(t, err)
		assert.Equal(t, big.NewInt(70), result.Amount())
		assert.Equal(t, "USD", result.Currency())
	})

	t.Run("should fail with different currencies", func(t *testing.T) {
		money1, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		money2, _ := valueobjects.NewMoney(big.NewInt(30), "EUR")

		_, err := money1.Subtract(money2)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot subtract different currencies")
	})

	t.Run("should fail when result would be negative", func(t *testing.T) {
		money1, _ := valueobjects.NewMoney(big.NewInt(30), "USD")
		money2, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		_, err := money1.Subtract(money2)

		assert.Error(t, err)
		assert.Equal(t, valueobjects.ErrNegativeBalance, err)
	})

	t.Run("subtracting zero should work", func(t *testing.T) {
		money1, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		money2, _ := valueobjects.NewMoney(big.NewInt(0), "USD")

		result, err := money1.Subtract(money2)

		require.NoError(t, err)
		assert.Equal(t, big.NewInt(100), result.Amount())
	})
}

func TestMoney_GreaterThanOrEqual(t *testing.T) {
	t.Run("should return true when amounts are equal", func(t *testing.T) {
		money1, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		money2, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		result, err := money1.GreaterThanOrEqual(money2)

		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when first amount is greater", func(t *testing.T) {
		money1, _ := valueobjects.NewMoney(big.NewInt(150), "USD")
		money2, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		result, err := money1.GreaterThanOrEqual(money2)

		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when first amount is less", func(t *testing.T) {
		money1, _ := valueobjects.NewMoney(big.NewInt(50), "USD")
		money2, _ := valueobjects.NewMoney(big.NewInt(100), "USD")

		result, err := money1.GreaterThanOrEqual(money2)

		require.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should fail with different currencies", func(t *testing.T) {
		money1, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		money2, _ := valueobjects.NewMoney(big.NewInt(100), "EUR")

		_, err := money1.GreaterThanOrEqual(money2)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot compare different currencies")
	})
}

func TestMoney_IsZero(t *testing.T) {
	t.Run("should return true for zero amount", func(t *testing.T) {
		money, _ := valueobjects.NewMoney(big.NewInt(0), "USD")
		assert.True(t, money.IsZero())
	})

	t.Run("should return false for non-zero amount", func(t *testing.T) {
		money, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		assert.False(t, money.IsZero())
	})
}

func TestMoney_IsNegative(t *testing.T) {
	t.Run("should return false for positive amount", func(t *testing.T) {
		money, _ := valueobjects.NewMoney(big.NewInt(100), "USD")
		assert.False(t, money.IsNegative())
	})

	t.Run("should return false for zero amount", func(t *testing.T) {
		money, _ := valueobjects.NewMoney(big.NewInt(0), "USD")
		assert.False(t, money.IsNegative())
	})
}
