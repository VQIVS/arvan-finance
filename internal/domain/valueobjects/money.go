package valueobjects

import (
	"errors"
	"math/big"
	"strings"
)

var (
	ErrNegativeBalance = errors.New("operation would result in negative balance")
)

type Money struct {
	amount   *big.Int
	currency string
}

func NewMoney(amount *big.Int, currency string) (Money, error) {
	if amount == nil {
		return Money{}, errors.New("amount cannot be nil")
	}
	if amount.Sign() < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}

	return Money{
		amount:   new(big.Int).Set(amount),
		currency: strings.ToUpper(currency),
	}, nil
}

func (m Money) GreaterThanOrEqual(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, errors.New("cannot compare different currencies")
	}
	return m.amount.Cmp(other.amount) >= 0, nil
}

func (m Money) IsNegative() bool {
	return m.amount.Sign() < 0
}

func (m Money) IsZero() bool {
	return m.amount.Sign() == 0
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("cannot subtract different currencies")
	}

	result := new(big.Int).Sub(m.amount, other.amount)

	if result.Sign() < 0 {
		return Money{}, ErrNegativeBalance
	}

	return Money{amount: result, currency: m.currency}, nil
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Amount() *big.Int {
	return new(big.Int).Set(m.amount)
}
