package routes

import (
	"example.com/rest_api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", canceRegistration)
	authenticated.GET("/events/:id/attendees", getAttendees)
	authenticated.GET("/registrations/me", getMyRegistrations)

	server.POST("/signup", signup)
	server.POST("/login", login)
}
