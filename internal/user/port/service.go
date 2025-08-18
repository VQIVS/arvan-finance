package port

import (
	"billing-service/internal/user/domain"
	"context"
)

type Service interface {
	CreateUser(ctx context.Context, user domain.User) (domain.APIKey, error)
	GetUserByID(ctx context.Context, ID domain.UserID) (domain.User, error)
	CreditUserBalance(ctx context.Context, ID domain.UserID, amount float64) error
	DebitUserBalance(ctx context.Context, ID domain.UserID, amount float64) error
}
