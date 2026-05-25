package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/vibe-be/pkg/tracing"
)

func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetHeader(tracing.HeaderName)
		if traceId == "" {
			traceId = tracing.NewTraceID()
		}

		c.Header(tracing.HeaderName, traceId)
		c.Set(tracing.TraceKey, traceId)
		c.Request = c.Request.WithContext(tracing.WithTraceID(c.Request.Context(), traceId))

		c.Next()
	}
}
