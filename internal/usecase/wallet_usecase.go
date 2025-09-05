package usecase

import (
	"context"
	"finance/internal/domain/entities"
	"finance/internal/domain/valueobjects"
	"math/big"

	"github.com/google/uuid"
)

type WalletService struct {
	WalletRepo entities.WalletRepo
	UserRepo   entities.UserRepo
	TxRepo     entities.TxRepo
}

func NewWalletService(walletRepo entities.WalletRepo, userRepo entities.UserRepo, txRepo entities.TxRepo) *WalletService {
	return &WalletService{
		WalletRepo: walletRepo,
		UserRepo:   userRepo,
		TxRepo:     txRepo,
	}
}

func (s *WalletService) DebitUserbalance(ctx context.Context, userID uuid.UUID, amount big.Int) error {
	wallet, err := s.WalletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	money, err := valueobjects.NewMoney(&amount, wallet.Currency)
	if err != nil {
		return err
	}
	if err := wallet.Debit(money); err != nil {
		return err
	}

	err = s.WalletRepo.UpdateBalance(ctx, wallet)
	if err != nil {
		return err
	}

	//TODO: move to up
	transaction := entities.NewTransaction(wallet.ID, userID, money, entities.TransactionDebit)
	err = s.TxRepo.Create(ctx, transaction)
	if err != nil {
		return err
	}
	transaction.MarkCompleted()
	err = s.TxRepo.UpdateStatus(ctx, transaction, entities.TransactionCompleted)
	if err != nil {
		return err
	}

	// add publisher here to publish tx completed event
	return nil
}

func (s *WalletService) CreditUserBalance(ctx context.Context, userID uuid.UUID, amount big.Int) error {
	wallet, err := s.WalletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	money, err := valueobjects.NewMoney(&amount, wallet.Currency)
	if err != nil {
		return err
	}
	if err := wallet.Credit(money); err != nil {
		return err
	}

	err = s.WalletRepo.UpdateBalance(ctx, wallet)
	if err != nil {
		return err
	}
	//TODO: move to up
	transaction := entities.NewTransaction(wallet.ID, userID, money, entities.TransactionCredit)
	err = s.TxRepo.Create(ctx, transaction)
	if err != nil {
		return err
	}
	transaction.MarkCompleted()
	err = s.TxRepo.UpdateStatus(ctx, transaction, entities.TransactionCompleted)
	if err != nil {
		return err
	}
	return nil
}

func (s *WalletService) RefundTransaction(ctx context.Context, txID string) error {
	tx, err := s.TxRepo.FindByID(ctx, txID)
	if err != nil {
		return err
	}

	if tx.Type != entities.TransactionDebit {
		return nil
	}

	wallet, err := s.WalletRepo.FindByID(ctx, tx.WalletID)
	if err != nil {
		return err
	}

	if err := wallet.Credit(tx.Amount); err != nil {
		return err
	}
	err = s.WalletRepo.UpdateBalance(ctx, wallet)
	if err != nil {
		return err
	}

	refundTx := entities.NewTransaction(wallet.ID, tx.UserID, tx.Amount, entities.TransactionCredit)
	err = s.TxRepo.Create(ctx, refundTx)
	if err != nil {
		return err
	}
	refundTx.MarkCompleted()
	err = s.TxRepo.UpdateStatus(ctx, refundTx, entities.TransactionCompleted)
	if err != nil {
		return err
	}
	// add publisher here to publish tx completed event
	return nil
}
