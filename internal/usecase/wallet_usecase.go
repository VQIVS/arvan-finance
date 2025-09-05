package usecase

import (
	"context"
	"finance/internal/domain/entities"
	"finance/internal/domain/valueobjects"
	"finance/internal/infra/storage"
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletService struct {
	WalletRepo      entities.WalletRepo
	UserRepo        entities.UserRepo
	TransactionRepo entities.TransactionRepo
	TxManager       storage.TransactionManager
}

func NewWalletService(walletRepo entities.WalletRepo, userRepo entities.UserRepo, txRepo entities.TransactionRepo, txManager storage.TransactionManager) *WalletService {
	return &WalletService{
		WalletRepo:      walletRepo,
		UserRepo:        userRepo,
		TransactionRepo: txRepo,
		TxManager:       txManager,
	}
}

func (s *WalletService) DebitUserbalance(ctx context.Context, userID uuid.UUID, amount big.Int) error {
	return s.withTransaction(func(walletRepo entities.WalletRepo, txRepo entities.TransactionRepo, userRepo entities.UserRepo) error {
		wallet, err := walletRepo.FindByUserID(ctx, userID)
		if err != nil {
			return err
		}

		money, err := valueobjects.NewMoney(&amount, wallet.Currency)
		if err != nil {
			return err
		}

		transaction := entities.NewTransaction(wallet.ID, userID, money, entities.TransactionDebit)
		if err := txRepo.Create(ctx, transaction); err != nil {
			return err
		}

		if err := wallet.Debit(money); err != nil {
			return err
		}

		if err := walletRepo.UpdateBalance(ctx, wallet); err != nil {
			return err
		}

		if err := transaction.MarkCompleted(); err != nil {
			return err
		}

		if err := txRepo.UpdateStatus(ctx, transaction, entities.TransactionCompleted); err != nil {
			return err
		}

		// add publisher here to publish tx completed event
		return nil
	})
}

func (s *WalletService) CreditUserBalance(ctx context.Context, userID uuid.UUID, amount big.Int) error {
	return s.withTransaction(func(walletRepo entities.WalletRepo, txRepo entities.TransactionRepo, userRepo entities.UserRepo) error {
		wallet, err := walletRepo.FindByUserID(ctx, userID)
		if err != nil {
			return err
		}

		money, err := valueobjects.NewMoney(&amount, wallet.Currency)
		if err != nil {
			return err
		}
		transaction := entities.NewTransaction(wallet.ID, userID, money, entities.TransactionCredit)
		if err := txRepo.Create(ctx, transaction); err != nil {
			return err
		}

		if err := wallet.Credit(money); err != nil {
			return err
		}

		if err := walletRepo.UpdateBalance(ctx, wallet); err != nil {
			return err
		}

		if err := transaction.MarkCompleted(); err != nil {
			return err
		}

		if err := txRepo.UpdateStatus(ctx, transaction, entities.TransactionCompleted); err != nil {
			return err
		}

		return nil
	})
}

func (s *WalletService) RefundTransaction(ctx context.Context, txID string) error {
	return s.withTransaction(func(walletRepo entities.WalletRepo, txRepo entities.TransactionRepo, userRepo entities.UserRepo) error {
		originalTx, err := txRepo.FindByID(ctx, txID)
		if err != nil {
			return err
		}

		if originalTx.Type != entities.TransactionDebit {
			return nil
		}

		wallet, err := walletRepo.FindByID(ctx, originalTx.WalletID)
		if err != nil {
			return err
		}

		refundTx := entities.NewTransaction(wallet.ID, originalTx.UserID, originalTx.Amount, entities.TransactionCredit)
		if err := txRepo.Create(ctx, refundTx); err != nil {
			return err
		}

		if err := wallet.Credit(originalTx.Amount); err != nil {
			return err
		}

		if err := walletRepo.UpdateBalance(ctx, wallet); err != nil {
			return err
		}

		if err := refundTx.MarkCompleted(); err != nil {
			return err
		}

		if err := txRepo.UpdateStatus(ctx, refundTx, entities.TransactionCompleted); err != nil {
			return err
		}

		// add publisher here to publish tx completed event
		return nil
	})
}

// withTransaction is a helper method that provides transactional repositories
func (s *WalletService) withTransaction(fn func(walletRepo entities.WalletRepo, txRepo entities.TransactionRepo, userRepo entities.UserRepo) error) error {
	return s.TxManager.WithTransaction(func(tx *gorm.DB) error {
		walletRepo := s.WalletRepo.WithTx(tx)
		txRepo := s.TransactionRepo.WithTx(tx)
		userRepo := s.UserRepo.WithTx(tx)
		return fn(walletRepo, txRepo, userRepo)
	})
}
