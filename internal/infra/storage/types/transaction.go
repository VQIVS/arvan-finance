package types

import "math/big"

type Money struct {
	Amount   *big.Int
	Currency string
}
type Transaction struct {
	Base
	WalletID string
	Wallet   *Wallet
	UserID   string
	User     *User
	Amount   Money
	Type     string
	Status   string
}
