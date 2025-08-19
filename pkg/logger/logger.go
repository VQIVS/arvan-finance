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
	return slog.New(slog.NewJSONHandler(os.Stdout, nil)).
		With("trace_id", uuid.NewString())
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
