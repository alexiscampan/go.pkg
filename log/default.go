package log

import (
	"context"
	l "log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultStore Store

func init() {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.EncoderConfig.MessageKey = "msg"
	config.EncoderConfig.TimeKey = "ts"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	defaultLogger, err := config.Build()
	if err != nil {
		l.Fatalln(err)
	}

	defaultStore = NewStore(defaultLogger)
}

// SetLogger defines the default package logger
func SetLogger(instance Store) {
	defaultStore = instance
}

// Bg delegates a no-context logger
func Bg() Logger {
	return defaultStore.Bg()
}

// For delegates a context logger
func For(ctx context.Context) Logger {
	return defaultStore.For(ctx)
}

// Default returns the logger factory
func Default() Store {
	return defaultStore
}
