package internalhttp

import (
	"time"

	"github.com/gin-gonic/gin"
)

func loggingMiddleware(c *gin.Context, app Application) {
	start := time.Now()

	c.Next()

	app.LogInfo(map[string]interface{}{
		"clientIP": c.ClientIP(),
		"time": time.Now().Format(time.RFC822),
		"method": c.Request.Method,
		"status": c.Writer.Status(),
		"latency": time.Since(start).Seconds(),
		"userAgent": c.Request.UserAgent(),
	})
}
