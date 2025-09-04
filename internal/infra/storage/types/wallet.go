package types

import "math/big"

type Wallet struct {
	Base
	UserID   string
	User     *User
	Balance  big.Int
	Currency string
}
