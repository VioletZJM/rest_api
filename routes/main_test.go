package routes

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"example.com/rest_api/database"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T) *gin.Engine {
	t.Helper()
	t.Setenv("DB_PATH", filepath.Join(t.TempDir(), "test.db"))
	t.Setenv("JWT_SECRET", "test")

	err := database.InitDB()
	if err != nil {
		t.Fatal(err)
	}

	e := gin.New()
	RegisterRoutes(e)
	return e
}

// 统一发送测试请求的函数
func doTest(t *testing.T, e *gin.Engine, method, path string, token string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func assertStatusCode(t *testing.T, rec *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if rec.Code != expected {
		t.Fatalf("Expected status code %d, but got %d,and body=%s", expected, rec.Code, rec.Body.String())
	}
}

func signupAndLogin(t *testing.T, e *gin.Engine, email string) string {
	rec := doTest(t, e, "POST", "/signup", "", `{"email": "`+email+`", "password": "123456"}`)
	assertStatusCode(t, rec, 201)

	rec = doTest(t, e, "POST", "/login", "", `{"email": "`+email+`", "password": "123456"}`)
	assertStatusCode(t, rec, 200)

	var resp struct {
		Token string `json:"token"`
	}

	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("token is empty,body=%s", rec.Body.String())
	}
	return resp.Token
}
