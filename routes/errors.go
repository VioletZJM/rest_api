package routes

import (
	"errors"
	"log"

	"example.com/rest_api/models"
	"example.com/rest_api/utils"
	"github.com/gin-gonic/gin"
)

func failWith(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrNotFound):
		utils.Fail(c, 404, "NOT_FOUND", "resource not found")
	case errors.Is(err, models.ErrAlreadyRegistered):
		utils.Fail(c, 409, "ALREADY_REGISTERED", "already registered")
	case errors.Is(err, models.ErrRegistrationNotFound):
		utils.Fail(c, 404, "REGISTRATION_NOT_FOUND", "registration not found")
	default:
		log.Printf("internal server error: %v", err)
		utils.Fail(c, 500, "INTERNAL_SERVER_ERROR", "internal server error")
	}
}
