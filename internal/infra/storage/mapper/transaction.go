package mapper

import (
	"finance/internal/domain/entities"
	"finance/internal/domain/valueobjects"
	"finance/internal/infra/storage/types"
)

// TODO: fix BASE
func moneyDomain2Storage(money valueobjects.Money) types.Money {
	return types.Money{
		Amount:   types.NewBigInt(money.Amount()),
		Currency: money.Currency(),
	}
}

func moneyStorage2Domain(money types.Money) (valueobjects.Money, error) {
	m, err := valueobjects.NewMoney(money.Amount.Int, money.Currency)
	if err != nil {
		return valueobjects.Money{}, err
	}
	return m, nil
}

func TxStorage2Domain(tx types.Transaction) (*entities.Transaction, error) {
	amount, err := moneyStorage2Domain(tx.Amount)
	if err != nil {
		return nil, err
	}
	return &entities.Transaction{
		ID:        tx.ID,
		WalletID:  tx.WalletID,
		UserID:    tx.UserID,
		Amount:    amount,
		Type:      entities.TransactionType(tx.Type),
		Status:    entities.TransactionStatus(tx.Status),
		SMSID:     tx.SMSID,
		CreatedAt: tx.CreatedAt,
		UpdatedAt: tx.UpdatedAt,
	}, nil
}

func TxDomain2Storage(tx *entities.Transaction) types.Transaction {
	return types.Transaction{
		Base:     types.Base{ID: tx.ID, CreatedAt: tx.CreatedAt, UpdatedAt: tx.UpdatedAt},
		WalletID: tx.WalletID,
		UserID:   tx.UserID,
		Amount:   moneyDomain2Storage(tx.Amount),
		Type:     string(tx.Type),
		Status:   string(tx.Status),
		SMSID:    tx.SMSID,
	}
}
