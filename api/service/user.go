package service

import (
	"billing-service/api/presenter"
	"billing-service/internal/user/domain"
	"billing-service/internal/user/port"
	"context"
)

type UserService struct {
	svc port.Service
}

func NewUserService(svc port.Service) *UserService {
	return &UserService{svc: svc}
}
func (s *UserService) CreateUser(ctx context.Context, req *presenter.UserRequest) (*presenter.UserResponse, error) {
	apiKey, err := s.svc.CreateUser(ctx, domain.User{
		APIKey:  req.APIKey,
		Balance: req.Balance,
	})
	if err != nil {
		return nil, err
	}
	return &presenter.UserResponse{
		APIKey: apiKey,
	}, nil
}

func (s *UserService) CreditUserBalance(ctx context.Context, req *presenter.UserBalanceRequest) (*presenter.UserBalanceResponse, error) {
	if err := s.svc.CreditUserBalance(ctx, req.UserID, req.Amount); err != nil {
		return nil, err
	}

	user, err := s.svc.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &presenter.UserBalanceResponse{
		UserID:  user.ID,
		Balance: user.Balance,
	}, nil

}

func (s *UserService) GetUserByID(ctx context.Context, userID domain.UserID) (*presenter.UserResponse, error) {
	user, err := s.svc.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &presenter.UserResponse{
		ID:      user.ID,
		APIKey:  user.APIKey,
		Balance: user.Balance,
	}, nil
}
