package internalhttp

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/domain"
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/pkg/apperror"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GetEventHandler(c *gin.Context, app Application) {
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	event, err := app.GetEventByID(context.Background(), uint(eventID))
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": apperror.ErrNotFound})
			return
		}
		if err = app.LogError(err); err != nil {
			log.Println(err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal})
		return
	}

	c.JSON(http.StatusOK, event)
}

func GetEventsHandler(c *gin.Context, app Application) {
	var params map[string]interface{}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	events, err := app.GetEventsByParams(context.Background(), params)
	if err != nil {
		if err = app.LogError(err); err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if len(events) == 0 {
		events = make([]domain.Event, 0)
	}

	c.JSON(http.StatusOK, events)
}

func UpdateEventHandler(c *gin.Context, app Application) {
	var event domain.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong params"})
		return
	}

	err = app.Update(context.Background(), event)
	if err != nil {
		if err = app.LogError(err); err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func CreateEventHandler(c *gin.Context, app Application) {
	log.Println("CreateEventHandler")
	var event domain.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong params"})
		return
	}

	newID, err := app.Create(context.Background(), event)
	if err != nil {
		if err = app.LogError(err); err != nil {
			log.Println(err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal})
		return
	}

	c.JSON(http.StatusOK, gin.H{"newID": newID})
}

func DeleteEventHandler(c *gin.Context, app Application) {
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong params"})
		return
	}

	err = app.Delete(context.Background(), uint(eventID))
	if err != nil {
		if err = app.LogError(err); err != nil {
			log.Println(err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.ErrInternal})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
