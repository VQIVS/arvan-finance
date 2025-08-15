package port

import (
	"billing-service/internal/transaction/domain"
	"context"
)

type Repo interface {
	Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	GetByFilter(ctx context.Context, filter domain.TransactionFilter) (*domain.Transaction, error)
}
