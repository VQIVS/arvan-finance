package storage

import (
	"billing-service/internal/user/domain"
	"billing-service/internal/user/port"
	"billing-service/pkg/adapters/storage/mapper"
	"billing-service/pkg/adapters/storage/types"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var user types.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", ID).First(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	user.Balance = amount

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
