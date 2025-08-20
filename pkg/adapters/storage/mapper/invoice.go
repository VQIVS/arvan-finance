package mapper

import (
	"billing-service/internal/invoice/domain"
	"billing-service/pkg/adapters/storage/types"

	"gorm.io/gorm"
)

func InvoiceDomain2Storage(invoice domain.Invoice) *types.Invoice {
	return &types.Invoice{
		Model:  gorm.Model{ID: invoice.ID, CreatedAt: invoice.CreatedAt, UpdatedAt: invoice.UpdatedAt},
		UserID: invoice.UserID,
		Amount: invoice.Amount,
		SMSID:  invoice.SMSID,
	}
}

func InvoiceStorage2Domain(invoice *types.Invoice) domain.Invoice {
	return domain.Invoice{
		ID:        invoice.ID,
		UserID:    invoice.UserID,
		Amount:    invoice.Amount,
		SMSID:     invoice.SMSID,
		CreatedAt: invoice.CreatedAt,
		UpdatedAt: invoice.UpdatedAt,
	}
}
