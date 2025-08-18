package http

import (
	"billing-service/app"
	"billing-service/config"
	"billing-service/docs"
	"fmt"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Run(appContainer app.App, cfg config.ServerConfig) error {
	router := fiber.New()

	api := router.Group("/api/v1", setUserContext)

	registerBillingAPI(appContainer, cfg, api)

	router.Use(cors.New())

	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.Schemes = []string{}
	docs.SwaggerInfo.BasePath = "/api/v1"

	router.Get("/swagger/*", adaptor.HTTPHandler(httpSwagger.Handler()))

	return router.Listen(fmt.Sprintf(":%d", cfg.HttpPort))
}

func registerBillingAPI(appContainer app.App, cfg config.ServerConfig, router fiber.Router) {
	userSvcGetter := userServiceGetter(appContainer, cfg)
	router.Put("/credit", CreditUserBalance(userSvcGetter))
	router.Post("/create", setTransaction(appContainer.DB()), CreateUser(userSvcGetter))
	router.Get("/:id", GetUserByID(userSvcGetter))
}
