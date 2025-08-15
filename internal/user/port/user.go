package port

import (
	"billing-service/internal/user/domain"
	"context"

	"gorm.io/gorm"
)

type Repo interface {
	WithTx(tx *gorm.DB) Repo
	Create(ctx context.Context, user domain.User) (domain.APIKey, error)
	GetByID(ctx context.Context, ID uint) (domain.User, error)
	CreditBalance(ctx context.Context, ID uint, amount int64) (domain.User, error)
	DebitBalance(ctx context.Context, ID uint, amount int64) (domain.User, error)
}
