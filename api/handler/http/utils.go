package http

import (
	"context"
)

// service factory pattern
type ServiceGetter[T any] func(context.Context) T
