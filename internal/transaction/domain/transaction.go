package domain

import "google.golang.org/genproto/googleapis/type/decimal"

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
	CreatedAt     string
}

type TransactionFilter struct {
	ID          TransactionID
	Type        TransactionType
	ReferenceID ReferenceID
	UserID      UserID
}
