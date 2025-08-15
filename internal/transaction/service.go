package transaction

import (
	"billing-service/internal/transaction/domain"
	"billing-service/internal/transaction/port"
	"context"
	"errors"
)

var (
	ErrTransactionOnCreate = errors.New("error on creating new transaction")
	ErrTransactionNotFound = errors.New("transaction not found")
)

type service struct {
	repo port.Repo
}

func NewService(repo port.Repo) port.Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	transaction, err := s.repo.Create(ctx, transaction)
	if err != nil {
		return nil, ErrTransactionOnCreate
	}
	return transaction, nil

}
func (s *service) GetByTransactionByFilter(ctx context.Context, filter domain.TransactionFilter) (*domain.Transaction, error) {
	transaction, err := s.repo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, ErrTransactionNotFound
	}

	return transaction, nil
}
