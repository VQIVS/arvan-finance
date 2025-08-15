package mapper

import (
	"billing-service/internal/user/domain"
	"billing-service/pkg/adapters/storage/types"

	"gorm.io/gorm"
)

func UserDoamin2Storage(userDomain domain.User) *types.User {
	return &types.User{
		Model: gorm.Model{
			ID:        uint(userDomain.ID),
			CreatedAt: userDomain.CreatedAt,
		},
		APIKey:  string(userDomain.APIKey),
		Balance: userDomain.Balance,
	}
}

func UserStorage2Domain(userStorage *types.User) *domain.User {
	return &domain.User{
		ID:        domain.UserID(userStorage.ID),
		APIKey:    domain.APIKey(userStorage.APIKey),
		Balance:   userStorage.Balance,
		CreatedAt: userStorage.CreatedAt,
	}
}
