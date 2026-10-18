package routes

import (
	"errors"
	"net/http"
	"strconv"

	"example.com/rest_api/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(c *gin.Context) {
	userId := c.GetInt64("userId")
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not get event_id."})
		return
	}

	registered, err := models.IsUserRegistered(eventId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register user for event."})
		return
	}
	if registered {
		c.JSON(http.StatusConflict, gin.H{"message": "User already registered for this event."})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event."})
		return
	}

	err = event.Register(userId)
	if errors.Is(err, models.ErrAlreadyRegistered) {
		c.JSON(http.StatusConflict, gin.H{"message": "User already registered for this event."})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register user for event."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registered!"})
}

func canceRegistration(c *gin.Context) {
	userId := c.GetInt64("userId")
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not get event_id."})
		return
	}

	var event models.Event
	event.ID = eventId

	err = event.CanceRegistration(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not cancel registration."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cancelled!"})
}
