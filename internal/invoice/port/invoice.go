package port

import (
	"billing-service/internal/invoice/domain"
	"context"
)

type Repo interface {
	Create(ctx context.Context, invoice *domain.Invoice) error
	GetListByUserID(ctx context.Context, userID string) ([]*domain.Invoice, error)
	Update(ctx context.Context, invoice *domain.Invoice) error
	GetByID(ctx context.Context, id string) (*domain.Invoice, error)
	Delete(ctx context.Context, id string) error
}
