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
// @Success      200      {string}  string                   "OK"
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
	return c.SendStatus(fiber.StatusOK)
}
