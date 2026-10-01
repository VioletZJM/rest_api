package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()

	RegisterRoutes(engine)

	req := httptest.NewRequest("POST", "/events", nil)
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("期望 401,实际状态码为:%d", rec.Code)
	}

}
