package entities

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRepo interface {
	GetByID(ctx context.Context, ID uuid.UUID) (*User, error)
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
