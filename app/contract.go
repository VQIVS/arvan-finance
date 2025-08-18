package app

import (
	"billing-service/config"
	userPort "billing-service/internal/user/port"
	"billing-service/pkg/adapters/rabbit"
	"context"

	"gorm.io/gorm"
)

type App interface {
	UserService(ctx context.Context) userPort.Service
	DB() *gorm.DB
	Config() config.Config
	Rabbit() *rabbit.Rabbit
}
