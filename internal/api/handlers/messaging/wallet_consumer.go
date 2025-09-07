package messaging

import (
	"context"
	"encoding/json"
	"finance/config"
	"finance/internal/domain/events"
	"finance/internal/usecase"
	"finance/pkg/logger"
	"finance/pkg/rabbit"
	"math/big"

	"github.com/google/uuid"
)

type ConsumerHandler struct {
	walletService *usecase.WalletService
	cfg           config.Config
	consumer      *rabbit.Consumer
	log           *logger.Logger
}

func NewConsumerHandler(walletService *usecase.WalletService, cfg config.Config, rabbitConn *rabbit.RabbitConn, logger *logger.Logger) *ConsumerHandler {
	return &ConsumerHandler{
		walletService: walletService,
		cfg:           cfg,
		consumer:      rabbit.NewConsumer(rabbitConn),
		log:           logger,
	}
}

func (h *ConsumerHandler) HandleDebitWallet(ctx context.Context, message []byte) error {
	var msg events.RequestSMSBilling
	err := json.Unmarshal(message, &msg)
	if err != nil {
		h.log.Error("Error unmarshaling message:", "error", err)
		return err
	}
	userID, err := uuid.Parse(msg.UserID)
	if err != nil {
		h.log.Error("Invalid user ID:", "error", err)
		return err
	}
	smsID, err := uuid.Parse(msg.SMSID)
	if err != nil {
		h.log.Error("Invalid SMS ID:", "error", err)
		return err
	}

	amount := big.NewInt(msg.Amount)

	test, err := h.walletService.DebitUserbalance(ctx, userID, smsID, *amount)
	if err != nil {
		h.log.Error("Error debiting user balance:", "error", err)
		return err
	}

	err = h.walletService.Publisher.PublishEvent(ctx, test)
	if err != nil {
		h.log.Error("Error publishing event:", "error", err)
		return err
	}

	h.log.Info(ctx, "Successfully debited user", "user_id", msg.UserID, "sms_id", msg.SMSID, "amount", msg.Amount)
	return nil
}

func (h *ConsumerHandler) HandleRefundTransaction(ctx context.Context, message []byte) error {
	var msg events.RequestBillingRefund
	err := json.Unmarshal(message, &msg)
	if err != nil {
		h.log.Error("Error unmarshaling message:", "error", err)
		return err
	}
	err = h.walletService.RefundTransaction(ctx, msg.TransactionID)
	if err != nil {
		h.log.Error("Error refunding transaction:", "error", err)
		return err
	}
	h.log.Info(ctx, "Successfully refunded transaction", "transaction_id", msg.TransactionID)
	return nil
}

func (h *ConsumerHandler) Run(ctx context.Context) error {
	if err := h.consumer.SetQos(1); err != nil {
		h.log.Error("Failed to set QoS", "error", err)
		return err
	}
	for _, queue := range h.cfg.RabbitMQ.Queues {
		switch queue.Name {
		case rabbit.DebitQueueName:
			h.consumer.Subscribe(queue.Name, func(message []byte) error {
				return h.HandleDebitWallet(ctx, message)
			})
		case rabbit.RefundQueueName:
			h.consumer.Subscribe(queue.Name, func(message []byte) error {
				return h.HandleRefundTransaction(ctx, message)
			})
		default:
			h.log.Logger.Warn("unknown queue in configuration", "queue", queue.Name)
		}
	}
	h.log.Logger.Info("starting SMS consumer")
	if err := h.consumer.StartConsume(); err != nil {
		h.log.Info(ctx, "failed to start consumer", "error", err)
		return err
	}

	<-ctx.Done()
	h.log.Logger.Info("SMS consumer stopped")
	return ctx.Err()
}
