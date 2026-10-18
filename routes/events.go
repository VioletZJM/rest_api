package routes

import (
	"net/http"
	"strconv"

	"example.com/rest_api/models"
	"github.com/gin-gonic/gin"
)

func getEvents(c *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch events.Try again."})
		return
	}
	c.JSON(http.StatusOK, events)
}

func createEvent(c *gin.Context) {

	var event models.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	userId := c.GetInt64("userId")
	event.UserID = userId

	err = event.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create event.Try again."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event created!", "event": event})
}

func getEvent(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cloud not get id."})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Cloud not fetch the event."})
		return
	}

	c.JSON(http.StatusOK, event)

}

func updateEvent(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cloud not get id."})
		return
	}

	userId := c.GetInt64("userId")

	event, err := models.GetEventByID(eventId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Cloud not fetch the event."})
		return
	}

	if event.UserID != userId {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update"})
		return
	}

	var updateEvent models.Event
	err = c.ShouldBindJSON(&updateEvent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	updateEvent.ID = eventId
	err = updateEvent.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cloud not update the event."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated succseefully"})

}

func deleteEvent(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cloud not get id."})
		return
	}

	userId := c.GetInt64("userId")

	event, err := models.GetEventByID(eventId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Cloud not fetch the event."})
		return
	}

	if event.UserID != userId {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to delete"})
		return
	}

	err = event.Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cloud not delete the event."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted succseefully"})

}

func getMyRegistrations(c *gin.Context) {
	registrations, err := models.GetRegistrationsForUser(c.GetInt64("userId"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch registrations.Try again."})
		return
	}
	c.JSON(http.StatusOK, registrations)
}

func getAttendees(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cloud not get id."})
		return
	}

	_, err = models.GetEventByID(eventId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Cloud not fetch the event."})
		return
	}

	attendees, err := models.GetAttendees(eventId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Cloud not fetch the attendees."})
		return
	}

	c.JSON(http.StatusOK, attendees)
}
