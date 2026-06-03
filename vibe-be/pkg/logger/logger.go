package logger

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vibe-be/pkg/tracing"
)

type Logger struct {
	z *zap.Logger
}

func New(level, env string) (*Logger, error) {
	lvl, err := zapcore.ParseLevel(level)
	if err != nil {
		lvl = zapcore.InfoLevel
	}

	var cfg zap.Config
	if env == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	cfg.Level = zap.NewAtomicLevelAt(lvl)

	z, err := cfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		return nil, err
	}
	return &Logger{z: z}, nil
}

// CtxLogger logs via zap AND mirrors each entry as an OpenTelemetry span event,
// so logs appear under the span's "Logs" tab in Jaeger.
type CtxLogger struct {
	z   *zap.Logger
	ctx context.Context
}

// Ctx returns a context-aware logger. When a recording span exists on ctx,
// every log entry is also added as a span event (Jaeger UI "Logs" tab) and
// Error-level entries mark the span as failed.
func (l *Logger) Ctx(ctx context.Context) *CtxLogger {
	z := l.z
	if id := tracing.FromContext(ctx); id != "" {
		z = z.With(zap.String(tracing.TraceKey, id))
	}
	return &CtxLogger{z: z, ctx: ctx}
}

func (c *CtxLogger) Debug(msg string, fields ...zap.Field) { c.emit(zapcore.DebugLevel, msg, fields) }
func (c *CtxLogger) Info(msg string, fields ...zap.Field)  { c.emit(zapcore.InfoLevel, msg, fields) }
func (c *CtxLogger) Warn(msg string, fields ...zap.Field)  { c.emit(zapcore.WarnLevel, msg, fields) }
func (c *CtxLogger) Error(msg string, fields ...zap.Field) { c.emit(zapcore.ErrorLevel, msg, fields) }
func (c *CtxLogger) Fatal(msg string, fields ...zap.Field) { c.emit(zapcore.FatalLevel, msg, fields) }

func (c *CtxLogger) emit(lvl zapcore.Level, msg string, fields []zap.Field) {
	if ce := c.z.Check(lvl, msg); ce != nil {
		ce.Write(fields...)
	}

	span := trace.SpanFromContext(c.ctx)
	if !span.IsRecording() {
		return
	}

	attrs := fieldsToAttrs(fields)
	attrs = append(attrs, attribute.String("log.severity", lvl.String()))
	span.AddEvent(msg, trace.WithAttributes(attrs...))

	if lvl >= zapcore.ErrorLevel {
		span.SetStatus(codes.Error, msg)
		for _, f := range fields {
			if f.Type == zapcore.ErrorType {
				if err, ok := f.Interface.(error); ok && err != nil {
					span.RecordError(err)
				}
			}
		}
	}
}

func fieldsToAttrs(fields []zap.Field) []attribute.KeyValue {
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}
	attrs := make([]attribute.KeyValue, 0, len(enc.Fields))
	for k, v := range enc.Fields {
		attrs = append(attrs, attribute.String(k, fmt.Sprintf("%v", v)))
	}
	return attrs
}

// Raw returns the underlying zap.Logger. Use only for interop
// (e.g. fxevent.ZapLogger). Prefer Ctx(ctx) for app code.
func (l *Logger) Raw() *zap.Logger { return l.z }

func (l *Logger) Named(name string) *Logger {
	return &Logger{z: l.z.Named(name)}
}

func (l *Logger) Sync() error { return l.z.Sync() }

// Convenience methods for logs outside any request (startup, lifecycle, background jobs).
func (l *Logger) Info(msg string, fields ...zap.Field)  { l.z.Info(msg, fields...) }
func (l *Logger) Warn(msg string, fields ...zap.Field)  { l.z.Warn(msg, fields...) }
func (l *Logger) Error(msg string, fields ...zap.Field) { l.z.Error(msg, fields...) }
func (l *Logger) Fatal(msg string, fields ...zap.Field) { l.z.Fatal(msg, fields...) }
func (l *Logger) Debug(msg string, fields ...zap.Field) { l.z.Debug(msg, fields...) }
