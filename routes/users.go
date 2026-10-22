package routes

import (
	"net/http"

	"example.com/rest_api/models"
	"example.com/rest_api/utils"
	"github.com/gin-gonic/gin"
)

func signup(c *gin.Context) {
	var user models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	err = user.Save()
	if err != nil {
		failWith(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})

}

func login(c *gin.Context) {
	var user models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", "invalid request body")
		return
	}

	err = user.ValidateCredentials()
	if err != nil {
		failWith(c, err)
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		failWith(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successfully", "token": token})
}
