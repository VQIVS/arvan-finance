package logger

import (
	"log/slog"
	"os"
	"sync"

	"github.com/google/uuid"
)

var (
	instance *slog.Logger
	once     sync.Once
)

func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func GetLogger() *slog.Logger {
	once.Do(func() {
		instance = NewLogger()
	})
	return instance
}

func SetLogger(logger *slog.Logger) {
	instance = logger
}

func WithTraceID(logger *slog.Logger, traceID string) *slog.Logger {
	return logger.With("trace_id", traceID)
}

func NewTracedLogger(traceID string) *slog.Logger {
	return GetLogger().With("trace_id", traceID)
}

func NewAutoTracedLogger() *slog.Logger {
	return GetLogger().With("trace_id", uuid.NewString())
}

func GetTracedLogger() *slog.Logger {
	return NewAutoTracedLogger()
}
