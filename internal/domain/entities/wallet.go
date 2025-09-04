package entities

import "finance/internal/domain/valueobjects"

type WalletRepo interface {
	Save(wallet *Wallet) error
	FindByID(id string) (*Wallet, error)
	FindByUserID(userID string) (*Wallet, error)
	UpdateBalance(wallet *Wallet) error
}

type Wallet struct {
	ID      string
	UserID  string
	Balance valueobjects.Money
}
