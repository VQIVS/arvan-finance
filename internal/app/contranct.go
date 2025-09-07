package app

import (
	"finance/config"
	"finance/internal/usecase"
	"finance/pkg/rabbit"

	"gorm.io/gorm"
)

type App interface {
	Config() config.Config
	DB() *gorm.DB
	RabbitConn() *rabbit.RabbitConn
	WalletService() *usecase.WalletService
}
