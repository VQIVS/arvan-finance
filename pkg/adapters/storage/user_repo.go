package storage

import (
	"billing-service/internal/user/domain"
	"billing-service/internal/user/port"
	"billing-service/pkg/adapters/storage/mapper"
	"billing-service/pkg/adapters/storage/types"
	"context"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) port.Repo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, user domain.User) (domain.APIKey, error) {
	u := mapper.UserDoamin2Storage(user)
	return domain.APIKey(u.APIKey), r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepo) GetByID(ctx context.Context, ID uint) (domain.User, error) {
	var user types.User
	if err := r.db.WithContext(ctx).Where("id = ?", ID).First(&user).Error; err != nil {
		return domain.User{}, err
	}
	return *mapper.UserStorage2Domain(&user), nil
}

func (r *userRepo) UpdateUserBalance(ctx context.Context, ID domain.UserID, amount float64) error {
	var user types.User
	if err := r.db.WithContext(ctx).Where("id = ?", ID).First(&user).Error; err != nil {
		return err
	}
	user.Balance = amount
	return r.db.WithContext(ctx).Save(&user).Error
}
