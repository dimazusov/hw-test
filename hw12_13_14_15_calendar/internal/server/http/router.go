package internalhttp

import (
	"github.com/gin-gonic/gin"
)

func NewGinRouter(app Application) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		loggingMiddleware(c, app)
		appMiddleware(c, app)
	})

	router.GET("/events", GetEventsHandler)
	router.GET("/event/:id", GetEventHandler)
	router.PUT("/event", UpdateEventHandler)
	router.POST("/event", CreateEventHandler)
	router.DELETE("/event/:id", DeleteEventHandler)

	return router
}
