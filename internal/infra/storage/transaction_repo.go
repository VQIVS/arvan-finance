package storage

import (
	"context"
	"finance/internal/domain/entities"
	"finance/internal/infra/storage/mapper"
	"finance/internal/infra/storage/types"

	"gorm.io/gorm"
)

type TransactionRepo struct {
	Db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) entities.TransactionRepo {
	return &TransactionRepo{
		Db: db,
	}
}

func (r *TransactionRepo) Create(ctx context.Context, tx *entities.Transaction) error {
	model := mapper.TxDomain2Storage(tx)
	return r.Db.WithContext(ctx).Create(&model).Error
}

func (r *TransactionRepo) FindByID(ctx context.Context, id string) (*entities.Transaction, error) {
	var model types.Transaction
	if err := r.Db.WithContext(ctx).Preload("Wallet").Preload("User").First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	tx, err := mapper.TxStorage2Domain(model)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *TransactionRepo) UpdateStatus(ctx context.Context, tx *entities.Transaction, status entities.TransactionStatus) error {
	tx.Status = status
	model := mapper.TxDomain2Storage(tx)
	return r.Db.WithContext(ctx).Model(&model).Update("status", status).Error
}

func (r *TransactionRepo) BeginDbTx() *gorm.DB {
	return r.Db.Begin()
}

func (r *TransactionRepo) WithTx(tx *gorm.DB) entities.TransactionRepo {
	return &TransactionRepo{
		Db: tx,
	}
}
