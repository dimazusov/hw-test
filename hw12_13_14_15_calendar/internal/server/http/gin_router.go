package internalhttp

import (
	"github.com/gin-gonic/gin"
)

func NewGinRouter(app Application) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		loggingMiddleware(c, app)
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}
