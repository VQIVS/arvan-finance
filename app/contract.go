package app

import (
	"billing-service/config"
	userPort "billing-service/internal/user/port"
	"context"

	"gorm.io/gorm"
)

type App interface {
	UserService(ctx context.Context) userPort.Service
	DB() *gorm.DB
	Config() config.Config
}
