package routes

import (
	"net/http"
	"strconv"

	"example.com/rest_api/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(c *gin.Context) {
	userId := c.GetInt64("userId")
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cloud not get event_id."})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cloud not fetch event."})
		return
	}

	err = event.Register(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cloud not regisetr user for event."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registered!"})
}

func canceRegistration(c *gin.Context) {
	userId := c.GetInt64("userId")
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cloud not get event_id."})
		return
	}

	var event models.Event
	event.ID = eventId

	err = event.CanceRegistration(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cloud not cance regisetration."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cancelled!"})
}
