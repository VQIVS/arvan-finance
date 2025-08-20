package domain

import "time"

type Status string
type Type string

const (
	InvoiceStatusPending   Status = "pending"
	InvoiceStatusCompleted Status = "completed"
	InvoiceStatusCancelled Status = "cancelled"
	InvoiceTypeCredit      Type   = "credit"
	InvoiceTypeDebit       Type   = "debit"
)

type Invoice struct {
	ID        uint
	UserID    uint
	Amount    float64
	SMSID     uint
	Status    string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
