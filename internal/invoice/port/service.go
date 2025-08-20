package port

import (
	"billing-service/internal/invoice/domain"
	"context"
)

type Service interface {
	CreateInvoice(ctx context.Context, invoice *domain.Invoice) error
	GetListByUserID(ctx context.Context, userID string) ([]*domain.Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *domain.Invoice) error
	GetInvoiceByID(ctx context.Context, ID string) (*domain.Invoice, error)
	DeleteInvoice(ctx context.Context, ID string) error
}
