package routes

import (
	"net/http"
	"strconv"

	"example.com/rest_api/models"
	"example.com/rest_api/utils"
	"github.com/gin-gonic/gin"
)

type eventListresponse struct {
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Items []models.Event `json:"Items"`
}

func getEvents(c *gin.Context) {
	page, limit := parsePagination(c)
	offset := (page - 1) * limit

	events, err := models.GetEvents(limit, offset)
	if err != nil {
		failWith(c, err)
		return
	}

	total, err := models.CountEvents()
	if err != nil {
		failWith(c, err)
		return
	}

	response := eventListresponse{
		Total: total,
		Page:  page,
		Limit: limit,
		Items: events,
	}

	c.JSON(http.StatusOK, response)
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
		utils.Fail(c, http.StatusForbidden, "FORBIDDEN", "not the owner of this event")
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
		utils.Fail(c, http.StatusForbidden, "FORBIDDEN", "not the owner of this event")
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
