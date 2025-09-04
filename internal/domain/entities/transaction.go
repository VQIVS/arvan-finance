package entities

import (
	"finance/internal/domain/valueobjects"
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	TransactionPending   TransactionStatus = "pending"
	TransactionCompleted TransactionStatus = "completed"
	TransactionFailed    TransactionStatus = "failed"
	TransactionCanceled  TransactionStatus = "canceled"
)

type TransactionType string

const (
	TransactionDebit  TransactionType = "debit"
	TransactionCredit TransactionType = "credit"
)

type TxRepo interface {
	Save(tx *Transaction) error
	FindByID(id string) (*Transaction, error)
	UpdateStatus(tx *Transaction, status TransactionStatus) error
}
type Transaction struct {
	ID          uuid.UUID          `json:"id"`
	WalletID    uuid.UUID          `json:"wallet_id"`
	Amount      valueobjects.Money `json:"amount"`
	Type        TransactionType    `json:"type"`
	Status      TransactionStatus  `json:"status"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
