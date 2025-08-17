package presenter

import "billing-service/internal/user/domain"

type UserRequest struct {
	APIKey  domain.APIKey
	Balance float64
}
type UserResponse struct {
	ID      domain.UserID `json:"id"`
	Balance float64       `json:"balance"`
	APIKey  domain.APIKey
}

type UserBalanceRequest struct {
	UserID domain.UserID `json:"user_id"`
	Amount float64       `json:"amount"`
}

type UserBalanceResponse struct {
	UserID  domain.UserID `json:"user_id"`
	Balance float64       `json:"balance"`
}
