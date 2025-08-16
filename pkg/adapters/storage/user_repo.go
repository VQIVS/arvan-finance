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
	return domain.APIKey(u.APIKey), r.db.Table("user").WithContext(ctx).Create(u).Error
}

func (r *userRepo) GetByID(ctx context.Context, ID uint) (domain.User, error) {
	var user types.User
	if err := r.db.Table("user").WithContext(ctx).Where("id = ?", ID).First(&user).Error; err != nil {
		return domain.User{}, err
	}
	return *mapper.UserStorage2Domain(&user), nil
}

func (r *userRepo) CreditBalance(ctx context.Context, ID uint, amount int64) (domain.User, error) {
	var user types.User
	if err := r.db.Table("user").WithContext(ctx).Where("id = ?", ID).First(&user).Error; err != nil {
		return domain.User{}, err
	}
	user.Balance += float64(amount)
	if err := r.db.Table("user").WithContext(ctx).Save(&user).Error; err != nil {
		return domain.User{}, err
	}
	return *mapper.UserStorage2Domain(&user), nil
}

func (r *userRepo) DebitBalance(ctx context.Context, ID uint, amount int64) (domain.User, error) {
	var user types.User
	if err := r.db.Table("user").WithContext(ctx).Where("id = ?", ID).First(&user).Error; err != nil {
		return domain.User{}, err
	}
	user.Balance -= float64(amount)
	if err := r.db.Table("user").WithContext(ctx).Save(&user).Error; err != nil {
		return domain.User{}, err
	}
	return *mapper.UserStorage2Domain(&user), nil
	// Implement HasSufficientBalance in Service Layer
}
