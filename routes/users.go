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
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = user.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})

}

func login(c *gin.Context) {
	var user models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = user.ValidateCredentials()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Could not authenticate user"})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successfully", "token": token})
}
