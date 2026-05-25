package logger

import (
	"context"

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

// Ctx returns a zap.Logger with trace_id from ctx attached (if any).
// Use this for any log emitted inside a request lifecycle.
func (l *Logger) Ctx(ctx context.Context) *zap.Logger {
	if id := tracing.FromContext(ctx); id != "" {
		return l.z.With(zap.String(tracing.TraceKey, id))
	}
	return l.z
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
