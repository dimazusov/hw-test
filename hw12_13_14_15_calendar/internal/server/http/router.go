package internalhttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func NewGinRouter(app Application) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		loggingMiddleware(c, app)
	})

	router.GET("/events", handle(app, GetEventsHandler))
	router.GET("/event/:id", handle(app, GetEventHandler))
	router.PUT("/event", handle(app, UpdateEventHandler))
	router.POST("/event", handle(app, CreateEventHandler))
	router.DELETE("/event/:id", handle(app, DeleteEventHandler))

	return router
}

func handle(app Application, handler func(c *gin.Context, app Application)) func(c *gin.Context) {
	fmt.Println("handle")
	return func(c *gin.Context) { handler(c, app) }
}
