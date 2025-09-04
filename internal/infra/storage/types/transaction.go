package types

import "math/big"

type Transaction struct {
	Base
	WalletID string
	Wallet   *Wallet
	UserID   string
	User     *User
	Amount   big.Int
	Type     string
	Status   string
}
