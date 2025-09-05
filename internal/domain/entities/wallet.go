package entities

import (
	"context"
	"errors"
	"finance/internal/domain/valueobjects"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInsufficientBalance = errors.New("insufficient wallet balance")
	ErrNegativeBalance     = errors.New("operation would result in negative balance")
	ErrInvalidAmount       = errors.New("amount must be positive")
	ErrWalletNotFound      = errors.New("wallet not found")
)

type WalletRepo interface {
	Save(ctx context.Context, wallet *Wallet) error
	FindByID(ctx context.Context, ID uuid.UUID) (*Wallet, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*Wallet, error)

	// uses db layer lock
	UpdateBalance(ctx context.Context, wallet *Wallet) error
	WithTx(tx *gorm.DB) WalletRepo
}

type Wallet struct {
	ID        uuid.UUID
	UserID    uuid.UUID
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
		ID:        uuid.New(),
		UserID:    userID,
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
	if w == nil {
		return ErrWalletNotFound
	}

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

func (w *Wallet) Credit(amount valueobjects.Money) error {
	if w == nil {
		return ErrWalletNotFound
	}

	if amount.IsZero() || amount.IsNegative() {
		return ErrInvalidAmount
	}

	if w.Balance.Currency() != amount.Currency() {
		return fmt.Errorf("currency mismatch: wallet has %s, crediting %s",
			w.Balance.Currency(), amount.Currency())
	}

	newBalance, err := w.Balance.Add(amount)
	if err != nil {
		return fmt.Errorf("credit calculation failed: %w", err)
	}

	w.Balance = newBalance
	w.UpdatedAt = time.Now()
	return nil
}
