package mapper

import (
	"finance/internal/domain/entities"
	"finance/internal/infra/storage/types"
)

func UserDomain2Storage(u *entities.User) types.User {
	return types.User{
		Base: types.Base{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Name:     u.Name,
		LastName: u.LastName,
		Phone:    u.Phone,
		WalletID: u.WalletID,
	}
}

func UserStorage2Domain(u types.User) *entities.User {
	return &entities.User{
		ID:        u.ID,
		Name:      u.Name,
		LastName:  u.LastName,
		Phone:     u.Phone,
		WalletID:  u.WalletID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
