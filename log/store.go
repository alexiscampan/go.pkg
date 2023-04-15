package log

import (
	"context"

	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

// Store define the default store of the logger
type Store interface {
	Bg() Logger
	For(context.Context) Logger
}

type store struct {
	logger *zap.Logger
}

// NewStore creates a new instance of a Store
func NewStore(logger *zap.Logger) Store {
	return &store{logger: logger}
}

// Bg implements Store
func (s store) Bg() Logger {
	return &logger{logger: s.logger}
}

// For implements Store
func (s store) For(ctx context.Context) Logger {
	if span := trace.FromContext(ctx); span != nil {
		return &spanLogger{span: span.SpanContext(), logger: s.logger}
	}
	return s.Bg()
}
