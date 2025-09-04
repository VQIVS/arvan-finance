package entities

import (
	"errors"
	"finance/internal/domain/valueobjects"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInsufficientBalance = errors.New("insufficient wallet balance")
	ErrNegativeBalance     = errors.New("operation would result in negative balance")
	ErrInvalidAmount       = errors.New("amount must be positive")
	ErrWalletNotFound      = errors.New("wallet not found")
)

type WalletRepo interface {
	Save(wallet *Wallet) error
	FindByID(id string) (*Wallet, error)
	FindByUserID(userID string) (*Wallet, error)
	UpdateBalance(wallet *Wallet) error
}

type Wallet struct {
	ID        string
	UserID    string
	Balance   valueobjects.Money
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWallet(userID uuid.UUID, currency string) (*Wallet, error) {
	zeroAmount, err := valueobjects.NewMoney(big.NewInt(0), currency)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		ID:        uuid.New().String(),
		UserID:    userID.String(),
		Balance:   zeroAmount,
		Currency:  currency,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (w *Wallet) HasSufficientBalance(amount valueobjects.Money) error {
	if w.Balance.Currency() != amount.Currency() {
		return fmt.Errorf("currency mismatch: wallet has %s, requested %s",
			w.Balance.Currency(), amount.Currency())
	}

	hasEnough, err := w.Balance.GreaterThanOrEqual(amount)
	if err != nil {
		return fmt.Errorf("balance comparison failed: %w", err)
	}

	if !hasEnough {
		return ErrInsufficientBalance
	}

	return nil
}

func (w *Wallet) Debit(amount valueobjects.Money) error {
	if amount.IsZero() || amount.IsNegative() {
		return ErrInvalidAmount
	}

	if err := w.HasSufficientBalance(amount); err != nil {
		return err
	}

	newBalance, err := w.Balance.Subtract(amount)
	if err != nil {
		return fmt.Errorf("debit calculation failed: %w", err)
	}

	if newBalance.IsNegative() {
		return ErrNegativeBalance
	}

	w.Balance = newBalance
	w.UpdatedAt = time.Now()
	return nil
}
