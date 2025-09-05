package storage

import (
	"context"
	"finance/internal/domain/entities"
	"finance/internal/infra/storage/mapper"
	"finance/internal/infra/storage/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository struct {
	Db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) entities.WalletRepo {
	return &WalletRepository{
		Db: db,
	}
}

func (r *WalletRepository) Save(ctx context.Context, wallet *entities.Wallet) error {
	model := mapper.WalletDomain2Storage(wallet)
	return r.Db.WithContext(ctx).Create(&model).Error
}

func (r *WalletRepository) FindByID(ctx context.Context, ID uuid.UUID) (*entities.Wallet, error) {
	var model types.Wallet
	if err := r.Db.WithContext(ctx).First(&model, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	res, err := mapper.WalletStorage2Domain(model)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *WalletRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*entities.Wallet, error) {
	var model types.Wallet
	if err := r.Db.WithContext(ctx).First(&model, "user_id = ?", userID.String()).Error; err != nil {
		return nil, err
	}
	res, err := mapper.WalletStorage2Domain(model)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *WalletRepository) UpdateBalance(ctx context.Context, wallet *entities.Wallet) error {
	model := mapper.WalletDomain2Storage(wallet)
	return r.Db.WithContext(ctx).Model(&model).Update("balance", model.Balance).Error
}
