package consumer

import (
	"billing-service/app"
	"billing-service/pkg/logger"
	"context"
)

type Handler struct {
	app app.App
}

func New(a app.App) *Handler {
	return &Handler{app: a}
}

func (h *Handler) Start(ctx context.Context) error {
	if h.app == nil || h.app.Rabbit() == nil {
		logger.NewLogger().Info("no rabbit configured, consumer won't start")
		return nil
	}
	svc := h.app.UserService(context.Background())
	if err := h.app.Rabbit().Consume("sms.user.balance.update", func(b []byte) error {
		sms, err := svc.DebitUserBalance(context.Background(), b)
		svc.UpdateSMSStatus(context.Background(), sms)
		if err != nil {
			logger.NewLogger().Error("failed to update sms status", "sms", sms, "error", err)
		}
		svc.UpdateSMSStatus(context.Background(), sms)
		return nil

	}); err != nil {
		return err
	}

	<-ctx.Done()
	h.app.Rabbit().Close()
	logger.NewLogger().Info("consumer stopped")
	return nil
}
