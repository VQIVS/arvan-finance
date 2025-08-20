package http

import (
	"billing-service/api/service"
	"billing-service/app"
	"billing-service/config"
	"context"
)

type userServiceProvider struct {
	appContainer app.App
	cfg          config.ServerConfig
}

func (p *userServiceProvider) GetUserService(ctx context.Context) *service.UserService {
	return service.NewUserService(p.appContainer.UserService(ctx))
}

func newUserServiceGetter(appContainer app.App, cfg config.ServerConfig) UserServiceGetter {
	return &userServiceProvider{
		appContainer: appContainer,
		cfg:          cfg,
	}
}
