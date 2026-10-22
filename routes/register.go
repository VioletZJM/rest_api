package routes

import (
	"net/http"
	"strconv"

	"example.com/rest_api/models"
	"example.com/rest_api/utils"
	"github.com/gin-gonic/gin"
)

func registerForEvent(c *gin.Context) {
	userId := c.GetInt64("userId")
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "INVALID_EVENT_ID", "invalid event id")
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		failWith(c, err)
		return
	}

	err = event.Register(userId)
	if err != nil {
		failWith(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registered!"})
}

func cancelRegistration(c *gin.Context) {
	userId := c.GetInt64("userId")
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "INVALID_EVENT_ID", "invalid event id")
		return
	}

	var event models.Event
	event.ID = eventId

	err = event.CancelRegistration(userId)
	if err != nil {
		failWith(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cancelled!"})
}
