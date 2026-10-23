package routes

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func parsePagination(c *gin.Context) (page, limit int) {
	page = atoiOrDefault(c.DefaultQuery("page", "1"), 1)
	limit = atoiOrDefault(c.DefaultQuery("limit", "20"), 20)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return page, limit
}

func atoiOrDefault(s string, def int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}
