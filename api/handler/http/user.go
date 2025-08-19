package http

import (
	"billing-service/api/presenter"
	"billing-service/internal/user/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// CreateUserHandler godoc
// @Summary Create user
// @Description Create a new user with an api key and initial balance
// @Tags users
// @Accept json
// @Produce json
// @Param user body presenter.UserRequest true "User request"
// @Success 200 {object} presenter.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /create [post]
func CreateUser(svcGetter UserServiceGetter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter.GetUserService(c.UserContext())
		var req presenter.UserRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		resp, err := svc.CreateUser(c.UserContext(), &req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(resp)
	}
}

// CreditUserBalanceHandler godoc
// @Summary Credit user balance
// @Description Add amount to a user's balance
// @Tags users
// @Accept json
// @Produce json
// @Param body body presenter.UserBalanceRequest true "Balance request"
// @Success 200 {object} presenter.UserBalanceResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /credit [put]
func CreditUserBalance(svcGetter UserServiceGetter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter.GetUserService(c.UserContext())
		var req presenter.UserBalanceRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		resp, err := svc.CreditUserBalance(c.UserContext(), &req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(resp)
	}
}

// GetUserByIDHandler godoc
// @Summary Get user by ID
// @Description Retrieve user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} presenter.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /{id} [get]
func GetUserByID(svcGetter UserServiceGetter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter.GetUserService(c.UserContext())
		idStr := c.Params("id")
		if idStr == "" {
			return fiber.ErrBadRequest
		}
		idUint, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return fiber.ErrBadRequest
		}
		resp, err := svc.GetUser(c.UserContext(), domain.UserID(idUint))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(resp)
	}

}
