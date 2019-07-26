package httputil

import (
	"fmt"

	"github.com/CzarSimon/text-service/go/pkg/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var requestLog = logger.GetDefaultLogger("httputil/requestLog")

// Logger request logging middleware.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		stop := createTimer()
		path := c.Request.URL.Path
		requestID := GetRequestID(c)
		requestLog.Info(fmt.Sprintf("Incomming request: %s %s", c.Request.Method, path), zap.String("requestId", requestID))

		c.Next()

		latency := fmt.Sprintf("%.2f ms", stop())
		requestLog.Info(fmt.Sprintf("Outgoing request: %s %s", c.Request.Method, path),
			zap.Int("status", c.Writer.Status()),
			zap.String("requestId", requestID),
			zap.String("latency", latency))
	}
}
