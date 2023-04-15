// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type spanLogger struct {
	logger *zap.Logger
	span   trace.SpanContext
}

func (l spanLogger) Debug(msg string, fields ...zapcore.Field) {
	fields = l.appendTraceID(fields...)
	l.logger.Debug(msg, fields...)
}

func (l spanLogger) Info(msg string, fields ...zapcore.Field) {
	fields = l.appendTraceID(fields...)
	l.logger.Info(msg, fields...)
}

func (l spanLogger) Error(msg string, fields ...zapcore.Field) {
	fields = l.appendTraceID(fields...)
	l.logger.Error(msg, fields...)
}

func (l spanLogger) Warn(msg string, fields ...zapcore.Field) {
	fields = l.appendTraceID(fields...)
	l.logger.Warn(msg, fields...)
}

func (l spanLogger) Fatal(msg string, fields ...zapcore.Field) {
	fields = l.appendTraceID(fields...)
	l.logger.Fatal(msg, fields...)
}

func (l spanLogger) Panic(msg string, fields ...zapcore.Field) {
	fields = l.appendTraceID(fields...)
	l.logger.Panic(msg, fields...)
}

func (l spanLogger) appendTraceID(fields ...zapcore.Field) []zapcore.Field {
	return append(fields, zap.String("TraceID", l.span.TraceID.String()))
}
