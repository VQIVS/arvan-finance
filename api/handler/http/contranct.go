package http

import (
	"billing-service/api/service"
	"context"
)

type UserServiceGetter interface {
	GetUserService(ctx context.Context) *service.UserService
}
