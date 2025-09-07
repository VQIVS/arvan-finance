package dto

type CreditWalletRequest struct {
	UserID string `json:"user_id" validate:"required,uuid4"`
	Amount int    `json:"amount" validate:"required,gt=0"`
}
