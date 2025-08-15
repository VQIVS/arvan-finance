package mapper

import (
	"billing-service/internal/transaction/domain"
	"billing-service/pkg/adapters/storage/types"

	"gorm.io/gorm"
)

func TNDomain2Storage(tn *domain.Transaction) *types.Transaction {
	return &types.Transaction{
		Model: gorm.Model{
			ID:        uint(tn.ID),
			CreatedAt: tn.CreatedAt,
		},
		UserID:        uint(tn.UserID),
		Amount:        tn.Amount,
		Type:          string(tn.Type),
		ReferenceID:   string(tn.ReferenceID),
		TransactionAt: tn.TransactionAt,
	}

}

func TNStorage2Domain(tn *types.Transaction) *domain.Transaction {
	return &domain.Transaction{
		ID:            tn.ID,
		UserID:        tn.UserID,
		Amount:        tn.Amount,
		Type:          domain.TransactionType(tn.Type),
		ReferenceID:   tn.ReferenceID,
		TransactionAt: tn.TransactionAt,
		CreatedAt:     tn.CreatedAt,
	}
}
