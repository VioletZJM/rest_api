package utils

import "github.com/gin-gonic/gin"

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Fail(c *gin.Context, status int, code string, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": apiError{code, message}})
}
