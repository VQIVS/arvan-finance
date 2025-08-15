package domain

import (
	"time"

	"google.golang.org/genproto/googleapis/type/decimal"
)

type (
	TransactionType string
	TransactionID   uint
	ReferenceID     string
	UserID          uint
)

const (
	TransactionTypeCredit TransactionType = "CREDIT"
	TransactionTypeDebit  TransactionType = "DEBIT"
)

type Transaction struct {
	ID            uint
	UserID        uint
	Amount        decimal.Decimal
	Type          TransactionType
	ReferenceID   string
	TransactionAt string
	CreatedAt     time.Time
}

type TransactionFilter struct {
	UserID      *uint
	Type        *TransactionType
	ReferenceID *string
	FromDate    *time.Time
	ToDate      *time.Time
	Limit       int
	Offset      int
}
