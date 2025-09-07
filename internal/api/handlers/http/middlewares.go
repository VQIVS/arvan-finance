package http

import (
	"finance/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func setTraceID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := logger.WithTraceID(c.Context())
		c.SetUserContext(ctx)

		traceID := logger.GetTraceID(ctx)
		c.Set("X-Trace-ID", traceID)
		return c.Next()
	}
}
