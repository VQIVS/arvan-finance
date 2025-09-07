package types

import (
	"github.com/google/uuid"
)

type Money struct {
	Amount   BigInt `gorm:"type:text;not null;default:'0'"`
	Currency string `gorm:"type:varchar(3);not null;default:'USD'"`
}

type Transaction struct {
	Base
	WalletID uuid.UUID `gorm:"type:uuid;index;not null"`
	UserID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Amount   Money     `gorm:"embedded;embeddedPrefix:amount_"`
	Type     string    `gorm:"type:varchar(10);index;not null"`
	Status   string    `gorm:"type:varchar(20);index;not null;default:'pending'"`
	SMSID    uuid.UUID `gorm:"type:uuid;index"`
}
