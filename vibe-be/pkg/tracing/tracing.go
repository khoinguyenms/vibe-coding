package tracing

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
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

// FromContext returns the active trace ID. If an OpenTelemetry span is
// recording on ctx, its 32-hex TraceID wins (so logs correlate with Jaeger).
// Otherwise it falls back to a value previously stored via WithTraceID.
func FromContext(ctx context.Context) string {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		return sc.TraceID().String()
	}
	if v, ok := ctx.Value(ctxKey{}).(string); ok {
		return v
	}
	return ""
}
