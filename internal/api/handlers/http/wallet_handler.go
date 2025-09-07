package http

import (
	"finance/internal/api/dto"
	"finance/internal/usecase"
	"math/big"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WalletHandler struct {
	walletService *usecase.WalletService
}

func NewWalletHandler(walletService *usecase.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// Credit godoc
// @Summary      Credit user wallet
// @Description  Credits a user's wallet with a specified amount
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreditWalletRequest  true  "Credit Wallet Request"
// @Success      200      {object}  dto.CreditWalletResponse "Wallet credited successfully"
// @Failure      400      {object}  map[string]interface{}   "Bad Request"
// @Failure      500      {object}  map[string]interface{}   "Internal Server Error"
// @Router       /wallet [post]
func (h *WalletHandler) Credit(c *fiber.Ctx) error {
	var req dto.CreditWalletRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user ID format")
	}

	amount := big.NewInt(int64(req.Amount))

	ctx := c.UserContext()
	if err := h.walletService.CreditUserBalance(ctx, userID, *amount); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(dto.BaseResponse{
		Success: true,
		Message: "Wallet credited successfully",
		Data: dto.CreditWalletResponse{
			UserID: userID.String(),
		},
	})
}

// GetWalletByUserID godoc
// @Summary      Get user wallet
// @Description  Gets a user's wallet information by user ID
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        user_id  path      string  true  "User ID"
// @Success      200      {object}  dto.GetWalletResponse "Wallet information"
// @Failure      400      {object}  map[string]interface{} "Bad Request"
// @Failure      404      {object}  map[string]interface{} "Wallet not found"
// @Failure      500      {object}  map[string]interface{} "Internal Server Error"
// @Router       /wallet/user/{user_id} [get]
func (h *WalletHandler) GetWalletByUserID(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	if userIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "user_id parameter is required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user ID format")
	}

	ctx := c.UserContext()
	wallet, err := h.walletService.GetWalletByUserID(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "wallet not found")
	}

	return c.JSON(dto.BaseResponse{
		Success: true,
		Message: "Wallet retrieved successfully",
		Data: dto.GetWalletResponse{
			ID:       wallet.ID.String(),
			UserID:   wallet.UserID.String(),
			Balance:  wallet.Balance.Amount().Int64(),
			Currency: wallet.Currency,
		},
	})
}

// GetUser godoc
// @Summary      Get user information
// @Description  Gets user information by user ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user_id  path      string  true  "User ID"
// @Success      200      {object}  dto.GetUserResponse "User information"
// @Failure      400      {object}  map[string]interface{} "Bad Request"
// @Failure      404      {object}  map[string]interface{} "User not found"
// @Failure      500      {object}  map[string]interface{} "Internal Server Error"
// @Router       /user/{user_id} [get]
func (h *WalletHandler) GetUser(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	if userIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "user_id parameter is required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user ID format")
	}

	ctx := c.UserContext()
	user, err := h.walletService.GetUserByID(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return c.JSON(dto.BaseResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data: dto.GetUserResponse{
			ID:       user.ID.String(),
			Name:     user.Name,
			LastName: user.LastName,
			Phone:    user.Phone,
			WalletID: user.WalletID,
		},
	})
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Gets a list of all users
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200      {object}  dto.GetAllUsersResponse "List of users"
// @Failure      500      {object}  map[string]interface{} "Internal Server Error"
// @Router       /users [get]
func (h *WalletHandler) GetAllUsers(c *fiber.Ctx) error {
	ctx := c.UserContext()
	users, err := h.walletService.GetAllUsers(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve users")
	}

	userResponses := make([]dto.GetUserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.GetUserResponse{
			ID:       user.ID.String(),
			Name:     user.Name,
			LastName: user.LastName,
			Phone:    user.Phone,
			WalletID: user.WalletID,
		}
	}

	return c.JSON(dto.BaseResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data: dto.GetAllUsersResponse{
			Users: userResponses,
			Total: len(userResponses),
		},
	})
}
