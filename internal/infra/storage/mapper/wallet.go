package mapper

import (
	"finance/internal/domain/entities"
	"finance/internal/domain/valueobjects"
	"finance/internal/infra/storage/types"
)

func WalletStorage2Domain(w types.Wallet) (*entities.Wallet, error) {
	money, err := valueobjects.NewMoney(w.Balance.Int, w.Currency)
	if err != nil {
		return nil, err
	}
	return &entities.Wallet{
		ID:        w.ID,
		UserID:    w.UserID,
		Balance:   money,
		Currency:  w.Currency,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}, nil
}

func WalletDomain2Storage(w *entities.Wallet) types.Wallet {
	return types.Wallet{
		Base: types.Base{
			ID:        w.ID,
			CreatedAt: w.CreatedAt,
			UpdatedAt: w.UpdatedAt,
		},
		UserID:   w.UserID,
		Balance:  types.NewBigInt(w.Balance.Amount()),
		Currency: w.Balance.Currency(),
	}
}
