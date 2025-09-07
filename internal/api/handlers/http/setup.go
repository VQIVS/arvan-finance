package http

import (
	"finance/config"
	"finance/internal/app"
	"fmt"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

func Run(appContainer app.App, cfg *config.Server) error {
	router := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.Schemes = []string{}
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.Get("/swagger/*", adaptor.HTTPHandler(httpSwagger.Handler()))

	return router.Listen(fmt.Sprintf(":%d", cfg.Port))
}
func registerSMSRoutes(appContainer app.App, router fiber.Router) {
	// ctx := context.Background()
	// walletUsecase := appContainer.WalletService(ctx)
	// walletHandler := NewwalletHandler(walletUsecase)

	v1 := router.Group("/api/v1")

	// SMS routes
	wallet := v1.Group("/wallet")
	wallet.Post("/", setTraceID())
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   "internal_error",
		"message": err.Error(),
		"code":    code,
	})
}
