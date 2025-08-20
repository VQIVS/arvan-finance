package invoice

import (
	"billing-service/internal/invoice/domain"
	"billing-service/internal/invoice/port"
	"context"
)

type service struct {
	repo port.Repo
}

func NewService(repo port.Repo) port.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateInvoice(ctx context.Context, invoice *domain.Invoice) error {
	return s.repo.Create(ctx, invoice)
}

func (s *service) GetListByUserID(ctx context.Context, userID string) ([]*domain.Invoice, error) {
	return s.repo.GetListByUserID(ctx, userID)
}

func (s *service) UpdateInvoice(ctx context.Context, invoice *domain.Invoice) error {
	return s.repo.Update(ctx, invoice)
}

func (s *service) GetInvoiceByID(ctx context.Context, ID string) (*domain.Invoice, error) {
	return s.repo.GetByID(ctx, ID)
}

func (s *service) DeleteInvoice(ctx context.Context, ID string) error {
	return s.repo.Delete(ctx, ID)
}
