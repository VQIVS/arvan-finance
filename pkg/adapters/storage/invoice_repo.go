package storage

import (
	"billing-service/internal/invoice/domain"
	"billing-service/internal/invoice/port"
	"billing-service/pkg/adapters/storage/mapper"
	"billing-service/pkg/adapters/storage/types"
	"context"

	"gorm.io/gorm"
)

type invoiceRepo struct {
	db *gorm.DB
}

func NewInvoiceRepo(db *gorm.DB) port.Repo {
	return &invoiceRepo{
		db: db,
	}
}

func (r *invoiceRepo) Create(ctx context.Context, invoice *domain.Invoice) error {
	storageInvoice := mapper.InvoiceDomain2Storage(*invoice)
	return r.db.Create(storageInvoice).Error
}

func (r *invoiceRepo) GetListByUserID(ctx context.Context, userID string) ([]*domain.Invoice, error) {
	var storageInvoices []*types.Invoice
	if err := r.db.Where("user_id = ?", userID).Find(&storageInvoices).Error; err != nil {
		return nil, err
	}

	var invoices []*domain.Invoice
	for _, storageInvoice := range storageInvoices {
		invoice := mapper.InvoiceStorage2Domain(storageInvoice)
		invoices = append(invoices, &invoice)
	}

	return invoices, nil
}

func (r *invoiceRepo) GetByID(ctx context.Context, id string) (*domain.Invoice, error) {
	var storageInvoice types.Invoice
	if err := r.db.First(&storageInvoice, "id = ?", id).Error; err != nil {
		return nil, err
	}

	invoice := mapper.InvoiceStorage2Domain(&storageInvoice)
	return &invoice, nil
}
func (r *invoiceRepo) Update(ctx context.Context, invoice *domain.Invoice) error {
	storageInvoice := mapper.InvoiceDomain2Storage(*invoice)
	return r.db.Save(storageInvoice).Error
}

func (r *invoiceRepo) Delete(ctx context.Context, id string) error {
	return r.db.Delete(&types.Invoice{}, "id = ?", id).Error
}
