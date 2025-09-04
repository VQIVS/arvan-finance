package entities

import (
	"errors"
	"finance/internal/domain/valueobjects"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionStatus string

const (
	TransactionPending   TransactionStatus = "pending"
	TransactionCompleted TransactionStatus = "completed"
	TransactionFailed    TransactionStatus = "failed"
)

var (
	ErrInvalidTransactionState = errors.New("invalid transaction state")
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

	//TODO: check the propper way to handle transactions in gorm
	BeginDbTx() *gorm.DB
}
type Transaction struct {
	ID        uuid.UUID          `json:"id"`
	WalletID  uuid.UUID          `json:"wallet_id"`
	UserID    uuid.UUID          `json:"user_id"`
	Amount    valueobjects.Money `json:"amount"`
	Type      TransactionType    `json:"type"`
	Status    TransactionStatus  `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func NewTransaction(walletID, userID uuid.UUID, amount valueobjects.Money, txType TransactionType) *Transaction {
	tx := &Transaction{
		ID:        uuid.New(),
		WalletID:  walletID,
		UserID:    userID,
		Amount:    amount,
		Type:      txType,
		Status:    TransactionPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return tx
}

func (t *Transaction) MarkCompleted() error {
	if t.Status != TransactionPending {
		return ErrInvalidTransactionState
	}
	t.Status = TransactionCompleted
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Transaction) MarkFailed() error {
	if t.Status != TransactionPending {
		return ErrInvalidTransactionState
	}
	t.Status = TransactionFailed
	t.UpdatedAt = time.Now()
	return nil
}
