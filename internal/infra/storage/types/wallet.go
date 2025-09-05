package types

import (
	"math/big"

	"github.com/google/uuid"
)

type Wallet struct {
	Base
	UserID   uuid.UUID
	User     *User
	Balance  big.Int
	Currency string
}
