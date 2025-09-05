package storage

import (
	"context"
	"finance/internal/domain/entities"
	"finance/internal/infra/storage/mapper"
	"finance/internal/infra/storage/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) entities.UserRepo {
	return &UserRepository{
		Db: db,
	}
}

func (r *UserRepository) GetByID(ctx context.Context, ID uuid.UUID) (*entities.User, error) {
	var model types.User
	if err := r.Db.WithContext(ctx).Preload("Wallet").First(&model, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	user := mapper.UserStorage2Domain(model)
	return user, nil
}
