package dto

type CreditWalletRequest struct {
	UserID string `json:"user_id" validate:"required,uuid4"`
	Amount int    `json:"amount" validate:"required,gt=0"`
}

type CreditWalletResponse struct {
	UserID string `json:"user_id"`
}

type GetWalletResponse struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type GetUserResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	LastName string  `json:"last_name"`
	Phone    string  `json:"phone"`
	WalletID *string `json:"wallet_id,omitempty"`
}

type GetAllUsersResponse struct {
	Users []GetUserResponse `json:"users"`
	Total int               `json:"total"`
}

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
