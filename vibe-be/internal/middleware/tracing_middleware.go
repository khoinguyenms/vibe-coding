package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/vibe-be/pkg/tracing"
)

// Tracing exposes the active trace ID on the response header and gin context.
// When otelgin runs ahead of this middleware (the normal case), the ID equals
// the OpenTelemetry trace ID, so X-Trace-Id and Jaeger traces line up.
func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := tracing.FromContext(c.Request.Context())
		if traceID == "" {
			traceID = c.GetHeader(tracing.HeaderName)
		}
		if traceID == "" {
			traceID = tracing.NewTraceID()
			c.Request = c.Request.WithContext(tracing.WithTraceID(c.Request.Context(), traceID))
		}

		c.Header(tracing.HeaderName, traceID)
		c.Set(tracing.TraceKey, traceID)

		c.Next()
	}
}
