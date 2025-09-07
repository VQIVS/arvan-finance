package types

import (
	"database/sql/driver"
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

type BigInt struct {
	*big.Int
}

func NewBigInt(value *big.Int) BigInt {
	if value == nil {
		return BigInt{big.NewInt(0)}
	}
	return BigInt{new(big.Int).Set(value)}
}

func (b *BigInt) Scan(value interface{}) error {
	if value == nil {
		b.Int = big.NewInt(0)
		return nil
	}

	switch v := value.(type) {
	case string:
		if b.Int == nil {
			b.Int = new(big.Int)
		}
		_, ok := b.Int.SetString(v, 10)
		if !ok {
			return fmt.Errorf("failed to parse big.Int from string: %s", v)
		}
	case []byte:
		if b.Int == nil {
			b.Int = new(big.Int)
		}
		_, ok := b.Int.SetString(string(v), 10)
		if !ok {
			return fmt.Errorf("failed to parse big.Int from bytes: %s", string(v))
		}
	case int64:
		if b.Int == nil {
			b.Int = new(big.Int)
		}
		b.Int.SetInt64(v)
	default:
		return fmt.Errorf("cannot scan %T into BigInt", value)
	}

	return nil
}

func (b BigInt) Value() (driver.Value, error) {
	if b.Int == nil {
		return "0", nil
	}
	return b.Int.String(), nil
}

type Wallet struct {
	Base
	UserID   uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	Balance  BigInt    `gorm:"type:text;not null;default:'0'"`
	Currency string    `gorm:"type:varchar(3);index;not null;default:'IRR'"`
}
