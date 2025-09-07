package dto

type CreditWalletRequest struct {
	UserID string `json:"user_id" validate:"required,uuid4"`
	Amount int    `json:"amount" validate:"required,gt=0"`
}

type CreditWalletResponse struct {
	UserID string `json:"user_id"`
}

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
