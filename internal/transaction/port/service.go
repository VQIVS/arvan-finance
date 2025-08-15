package port

import (
	"billing-service/internal/transaction/domain"
	"context"
)

type Service interface {
	CreateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	GetTransactionsByFilter(ctx context.Context, filter *domain.TransactionFilter) ([]*domain.Transaction, error)
}
