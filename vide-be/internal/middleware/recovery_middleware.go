package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/vibe-be/pkg/logger"
	"go.uber.org/zap"
)

func Recovery(log *logger.Logger) gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, recovered any) {
		stack := debug.Stack()

		log.Ctx(c.Request.Context()).Error("panic recovered",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("remote", c.ClientIP()),
			zap.Any("error", recovered),
			zap.ByteString("stack", stack))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
	})
}
