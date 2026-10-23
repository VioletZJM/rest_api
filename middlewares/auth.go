package middlewares

import (
	"strings"

	"example.com/rest_api/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		utils.Fail(c, 401, "UNAUTHORIZED", "missing or invalid token")
		return
	}

	authHeader, ok := strings.CutPrefix(token, "Bearer ")
	if !ok {
		utils.Fail(c, 401, "UNAUTHORIZED", "missing or invalid token")
		return
	}

	token = strings.TrimSpace(authHeader)
	if token == "" {
		utils.Fail(c, 401, "UNAUTHORIZED", "missing or invalid token")
		return
	}

	userId, err := utils.VerifyToken(token)
	if err != nil {
		utils.Fail(c, 401, "UNAUTHORIZED", "missing or invalid token")
		return
	}

	c.Set("userId", userId)

	c.Next()
}
