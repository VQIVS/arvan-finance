package types

import (
	"google.golang.org/genproto/googleapis/type/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID        uint
	Amount        decimal.Decimal
	Type          string
	ReferenceID   string
	TransactionAt string // business perspective
}
