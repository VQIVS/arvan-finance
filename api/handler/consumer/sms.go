package consumer

import (
	"billing-service/app"
	"billing-service/pkg/adapters/rabbit"
	"billing-service/pkg/constants"
	"context"

	"github.com/google/uuid"
)

type Handler struct {
	app app.App
}

func New(a app.App) *Handler {
	return &Handler{app: a}
}

func (h *Handler) Start(ctx context.Context) error {
	if h.app == nil || h.app.Rabbit() == nil {
		h.app.Logger().With("trace_id", uuid.NewString()).Info("no rabbit configured, consumer won't start")
		return nil
	}

	queue := rabbit.GetQueueName(constants.KeyBalanceUpdate)
	svc := h.app.UserService(context.Background())

	if err := h.app.Rabbit().Consume(queue, func(b []byte) error {
		sms, err := svc.DebitUserBalance(context.Background(), b)
		if updateErr := svc.UpdateSMSStatus(context.Background(), sms); updateErr != nil {
			h.app.Logger().With("trace_id", uuid.NewString()).Error("failed to update sms status", "sms", sms, "error", updateErr)
		}
		if err != nil {
			h.app.Logger().With("trace_id", uuid.NewString()).Error("failed to update sms status", "sms", sms, "error", err)
		}
		return nil

	}); err != nil {
		return err
	}

	<-ctx.Done()
	h.app.Logger().With("trace_id", uuid.NewString()).Info("consumer stopped")
	return nil
}
