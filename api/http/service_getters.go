package http

import (
	"billing-service/api/service"
	"billing-service/app"
	"billing-service/config"
	"context"
)

// dep injection pattern
// user service transient instance handler
func userServiceGetter(appContainer app.App, cfg config.ServerConfig) ServiceGetter[*service.UserService] {
	return func(ctx context.Context) *service.UserService {
		return service.NewUserService(appContainer.UserService(ctx))
	}
}
