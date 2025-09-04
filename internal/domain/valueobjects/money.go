package valueobjects

import (
	"errors"
	"math/big"
	"strings"
)

type Money struct {
	amount   *big.Int
	currency string
}

func NewMoney(amount *big.Int, currency string) (Money, error) {
	if amount == nil {
		return Money{}, errors.New("amount cannot be nil")
	}
	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}
	currency = strings.ToUpper(currency)

	if len(currency) != 3 {
		return Money{}, errors.New("currency must be 3-letter code")
	}
	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

func (m Money) Amount() *big.Int {
	return new(big.Int).Set(m.amount)
}

func (m Money) Currency() string {
	return m.currency
}
