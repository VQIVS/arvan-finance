package http

import (
	"billing-service/pkg/context"
	"billing-service/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func setUserContext(c *fiber.Ctx) error {
	c.SetUserContext(context.NewAppContext(c.UserContext(), context.WithLogger(logger.NewLogger())))
	return c.Next()
}

func setTransaction(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tx := db.Begin()

		context.SetDB(c.UserContext(), tx, true)

		err := c.Next()

		if c.Response().StatusCode() >= 300 {
			return context.Rollback(c.UserContext())
		}

		if err := context.CommitOrRollback(c.UserContext(), true); err != nil {
			return err
		}

		return err
	}
}
