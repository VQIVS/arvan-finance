package port

import (
	"billing-service/internal/transaction/domain"
	"context"

	"gorm.io/gorm"
)

type Repo interface {
	WithTx(tx *gorm.DB) Repo
	Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	GetByFilter(ctx context.Context, filter *domain.TransactionFilter) ([]*domain.Transaction, error)
}
