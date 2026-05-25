package tracing

import (
	"context"

	"github.com/google/uuid"
)

const (
	HeaderName = "X-Trace-Id"
	TraceKey   = "trace_id"
)

type ctxKey struct{}

func NewTraceID() string {
	return uuid.NewString()
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ctxKey{}, traceID)
}

func FromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKey{}).(string); ok {
		return v
	}
	return ""
}
