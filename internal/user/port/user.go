package port

import (
	"billing-service/internal/user/domain"
	"context"
)

type Repo interface {
	Create(ctx context.Context, user domain.User) (domain.APIKey, error)
	GetByID(ctx context.Context, ID uint) (domain.User, error)
	UpdateUserBalance(ctx context.Context, ID domain.UserID, amount float64) error
}
