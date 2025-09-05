package storage

import (
	"gorm.io/gorm"
)

type WalletRepository struct {
	Db *gorm.DB
}

// func NewWalletRepository(db *gorm.DB) entities.WalletRepo {
// 	return &WalletRepository{
// 		Db: db,
// 	}
// }
