package user

import (
	"billing-service/internal/user/domain"
	"billing-service/internal/user/port"
	"context"
	"errors"
)

var (
	ErrUserOnCreate        = errors.New("error on creating new user")
	ErrUserNotFound        = errors.New("user not found")
	ErrInsufficientBalance = errors.New("insufficient balance for debit operation")
)

type service struct {
	repo port.Repo
}

func NewService(repo port.Repo) port.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUser(ctx context.Context, user domain.User) (domain.APIKey, error) {
	apiKey, err := s.repo.Create(ctx, user)
	if err != nil {
		return "", ErrUserOnCreate
	}
	return apiKey, nil
}

func (s *service) GetUserByID(ctx context.Context, ID domain.UserID) (domain.User, error) {
	user, err := s.repo.GetByID(ctx, uint(ID))
	if err != nil {
		return domain.User{}, ErrUserNotFound
	}
	return user, nil
}

func (s *service) CreditUserBalance(ctx context.Context, ID domain.UserID, amount int64) error {
	_, err := s.repo.CreditBalance(ctx, uint(ID), amount)
	if err != nil {
		return ErrUserNotFound
	}
	return nil
}
func (s *service) DebitUserBalance(ctx context.Context, ID domain.UserID, amount int64) error {
	user, err := s.repo.DebitBalance(ctx, uint(ID), amount)
	if err != nil {
		return ErrUserNotFound
	}
	if !user.HasSufficientBalance(uint64(amount)) {
		return ErrInsufficientBalance
	}
	return nil
}
