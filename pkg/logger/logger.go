package logger

import (
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Logger struct{ sugar *zap.SugaredLogger }

var (
	defaultSugar *zap.SugaredLogger
	defaultOnce  sync.Once
)

func NewLogger() *Logger {
	defaultOnce.Do(func() {
		l, err := zap.NewProduction()
		if err != nil {
			l = zap.NewExample()
		}
		defaultSugar = l.Sugar()
	})
	return &Logger{sugar: defaultSugar.With("trace_id", uuid.NewString())}
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) > 0 {
		l.sugar.Infow(msg, keysAndValues...)
		return
	}
	l.sugar.Info(msg)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) > 0 {
		l.sugar.Errorw(msg, keysAndValues...)
		return
	}
	l.sugar.Error(msg)
}

func (l *Logger) With(keysAndValues ...interface{}) *Logger {
	return &Logger{sugar: l.sugar.With(keysAndValues...)}
}
