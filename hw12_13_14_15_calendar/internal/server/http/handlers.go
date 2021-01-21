package internalhttp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/domain"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/pkg/apperror"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GetEventHandler(c *gin.Context) {
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong params"})
		return
	}

	app := getApp(c)

	event, err := app.GetEventByID(context.Background(), uint(eventID))
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": apperror.ErrNotFound})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": apperror.ErrInternal})
	}

	c.JSON(http.StatusOK, event)
}

func GetEventsHandler(c *gin.Context) {
	var params map[string]interface{}

	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": apperror.ErrInternal})
		return
	}

	app := getApp(c)

	event, err := app.GetEventsByParams(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": apperror.ErrInternal})
	}

	c.JSON(http.StatusOK, event)
}

func UpdateEventHandler(c *gin.Context) {
	var event domain.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong params"})
		return
	}

	app := getApp(c)

	err = app.Update(context.Background(), event)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": apperror.ErrInternal})
	}

	c.JSON(http.StatusOK, gin.H{})
}

func CreateEventHandler(c *gin.Context) {
	var event domain.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong params"})
		return
	}

	app := getApp(c)

	events, err := app.Create(context.Background(), event)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": apperror.ErrInternal})
	}

	c.JSON(http.StatusOK, events)
}

func DeleteEventHandler(c *gin.Context) {
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong params"})
		return
	}

	app := getApp(c)

	err = app.Delete(context.Background(), uint(eventID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": apperror.ErrInternal})
	}

	c.JSON(http.StatusOK, gin.H{})
}

func getApp(c *gin.Context) Application {
	app, _ := c.Get("app")

	return app.(Application)
}
