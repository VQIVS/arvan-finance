package entities

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo interface {
	GetByID(ctx context.Context, ID uuid.UUID) (*User, error)
	WithTx(tx *gorm.DB) UserRepo
}

type User struct {
	ID        uuid.UUID
	Name      string
	LastName  string
	Phone     string
	WalletID  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
