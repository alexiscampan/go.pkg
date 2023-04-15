// Package log to define custom logger
package log

import (
	"fmt"

	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"
)

func init() {
	fmt.Println("logger init")
}

// Logger defines all the available levels
type Logger interface {
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	Panic(msg string, fields ...zapcore.Field)
}

type logger struct {
	logger *zap.Logger
}

func (l logger) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

func (l logger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

func (l logger) Warn(msg string, fields ...zapcore.Field) {
	l.logger.Warn(msg, fields...)
}

func (l logger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

func (l logger) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l logger) Panic(msg string, fields ...zapcore.Field) {
	l.logger.Panic(msg, fields...)
}
