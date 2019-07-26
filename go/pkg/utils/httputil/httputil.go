package httputil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthFunc health check function signature.
type HealthFunc func() error

// NewRouter creates a default router.
func NewRouter(healthCheck HealthFunc) *gin.Engine {
	r := gin.New()
	r.Use(
		gin.Recovery(),
		Metrics(),
		RequestID(),
		Locale(),
		Logger(),
		HandleErrors())

	r.GET("/health", checkHealth(healthCheck))
	r.GET(metricsPath, prometheusHandler())
	return r
}

// SendOK sends an ok status and message to the client.
func SendOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func checkHealth(check HealthFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := check()
		if err == nil {
			SendOK(c)
			return
		}

		c.Error(err)
	}
}
