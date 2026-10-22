package routes

import (
	"net/http"
	"strconv"

	"example.com/rest_api/models"
	"example.com/rest_api/utils"
	"github.com/gin-gonic/gin"
)

func getEvents(c *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		failWith(c, err)
		return
	}
	c.JSON(http.StatusOK, events)
}

func createEvent(c *gin.Context) {

	var event models.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	userId := c.GetInt64("userId")
	event.UserID = userId

	err = event.Save()
	if err != nil {
		failWith(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event created!", "event": event})
}

func getEvent(c *gin.Context) {
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

	c.JSON(http.StatusOK, event)

}

func updateEvent(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "INVALID_EVENT_ID", "invalid event id")
		return
	}

	userId := c.GetInt64("userId")

	event, err := models.GetEventByID(eventId)
	if err != nil {
		failWith(c, err)
		return
	}

	if event.UserID != userId {
		utils.Fail(c, http.StatusForbidden, "forbidden", "not the owner of this event")
		return
	}

	var updateEvent models.Event
	err = c.ShouldBindJSON(&updateEvent)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	updateEvent.ID = eventId
	err = updateEvent.Update()
	if err != nil {
		failWith(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})

}

func deleteEvent(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "INVALID_EVENT_ID", "invalid event id")
		return
	}

	userId := c.GetInt64("userId")

	event, err := models.GetEventByID(eventId)
	if err != nil {
		failWith(c, err)
		return
	}

	if event.UserID != userId {
		utils.Fail(c, http.StatusForbidden, "forbidden", "not the owner of this event")
		return
	}

	err = event.Delete()
	if err != nil {
		failWith(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})

}

func getMyRegistrations(c *gin.Context) {
	registrations, err := models.GetRegistrationsForUser(c.GetInt64("userId"))
	if err != nil {
		failWith(c, err)
		return
	}
	c.JSON(http.StatusOK, registrations)
}

func getAttendees(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "INVALID_EVENT_ID", "invalid event id")
		return
	}

	_, err = models.GetEventByID(eventId)
	if err != nil {
		failWith(c, err)
		return
	}

	attendees, err := models.GetAttendees(eventId)
	if err != nil {
		failWith(c, err)
		return
	}

	c.JSON(http.StatusOK, attendees)
}
