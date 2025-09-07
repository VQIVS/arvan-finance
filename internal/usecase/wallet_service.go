package usecase

import (
	"context"
	"finance/internal/domain/entities"
	"finance/internal/domain/events"
	"finance/internal/domain/valueobjects"
	"finance/internal/infra/storage"
	"finance/pkg/logger"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletService struct {
	WalletRepo      entities.WalletRepo
	UserRepo        entities.UserRepo
	TransactionRepo entities.TransactionRepo
	TxManager       storage.TransactionManager
	Publisher       events.Publisher
	log             *logger.Logger
}

func NewWalletService(walletRepo entities.WalletRepo,
	userRepo entities.UserRepo,
	transactionRepo entities.TransactionRepo,
	txManager storage.TransactionManager,
	publisher events.Publisher, log *logger.Logger) *WalletService {
	return &WalletService{
		WalletRepo:      walletRepo,
		UserRepo:        userRepo,
		TransactionRepo: transactionRepo,
		TxManager:       txManager,
		Publisher:       publisher,
		log:             log,
	}
}

// consumer handler calls this usecase
func (s *WalletService) DebitUserbalance(ctx context.Context, userID, smsID uuid.UUID, amount big.Int) error {
	var eventToPublish *events.SMSDebited

	err := s.withTransaction(func(walletRepo entities.WalletRepo, txRepo entities.TransactionRepo, userRepo entities.UserRepo) error {
		wallet, err := walletRepo.FindByUserID(ctx, userID)
		if err != nil {
			return err
		}

		money, err := valueobjects.NewMoney(&amount, wallet.Currency)
		if err != nil {
			return err
		}

		transaction := entities.NewTransaction(wallet.ID, userID, smsID, money, entities.TransactionDebit)
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

		// Prepare event for publishing after transaction commits
		eventToPublish = &events.SMSDebited{
			UserID:        userID.String(),
			SMSID:         transaction.SMSID.String(),
			Amount:        transaction.Amount.Amount().Int64(),
			TransactionID: transaction.ID.String(),
			TimeStamp:     time.Now(),
		}
		return nil
	})

	if err == nil && eventToPublish != nil {
		if publishErr := s.Publisher.PublishEvent(ctx, eventToPublish); publishErr != nil {
		}
		s.log.Info(ctx, "Published SMSDebited event", "user_id", eventToPublish.UserID, "sms_id", eventToPublish.SMSID, "amount", eventToPublish.Amount, "transaction_id", eventToPublish.TransactionID)
	}
	return err
}

// http handler calls this usecase
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
		transaction := entities.NewTransaction(wallet.ID, userID, uuid.New(), money, entities.TransactionCredit)
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

// this usecase executes in a subsciber handler
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

		refundTx := entities.NewTransaction(wallet.ID, originalTx.UserID, originalTx.SMSID, originalTx.Amount, entities.TransactionCredit)
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
