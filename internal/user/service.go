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
		return "", err
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

func (s *service) CreditUserBalance(ctx context.Context, ID domain.UserID, amount float64) error {
	if amount <= 0 {
		return errors.New("invalid credit amount")
	}
	user, err := s.repo.GetByID(ctx, uint(ID))
	if err != nil {
		return ErrUserNotFound
	}

	user.Balance += amount

	err = s.repo.UpdateUserBalance(ctx, ID, user.Balance)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DebitUserBalance(ctx context.Context, ID domain.UserID, amount float64) error {
	if amount <= 0 {
		return errors.New("invalid debit amount")
	}
	user, err := s.repo.GetByID(ctx, uint(ID))
	if err != nil {
		return ErrUserNotFound
	}
	if user.Balance < amount {
		return ErrInsufficientBalance
	}
	err = s.repo.UpdateUserBalance(ctx, ID, user.Balance-amount)
	if err != nil {
		return err
	}
	return nil
}
