package types

import (
	"math/big"

	"github.com/google/uuid"
)

type Money struct {
	Amount   *big.Int
	Currency string
}
type Transaction struct {
	Base
	WalletID uuid.UUID
	Wallet   *Wallet
	UserID   uuid.UUID
	User     *User
	Amount   Money
	Type     string
	Status   string
	SMSID    uuid.UUID
}
