package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options defines the options available for the logger
type Options struct {
	Debug    bool
	LogLevel string
}

// DefaultOptions define default options
var DefaultOptions = &Options{
	Debug:    true,
	LogLevel: "info",
}

// Setup our logger
func Setup(ctx context.Context, opts *Options) {
	config := zap.NewProductionConfig()

	config.Development = true
	config.Sampling = nil
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "msg"
	config.EncoderConfig.TimeKey = "ts"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true

	err := config.Level.UnmarshalText([]byte(opts.LogLevel))
	if err != nil {
		panic(err)
	}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	logStore := NewStore(logger)

	SetLogger(logStore)

	zap.ReplaceGlobals(logger)
}
