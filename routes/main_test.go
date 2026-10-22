package routes

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"example.com/rest_api/database"
	"example.com/rest_api/utils"
	"github.com/gin-gonic/gin"
)

const (
	testPassword   = "123456"
	validEventBody = `{"name":"Test Event","description":"desc","location":"online","dateTime":"2026-01-01T10:00:00Z"}`
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T) *gin.Engine {
	t.Helper()
	t.Setenv("DB_PATH", filepath.Join(t.TempDir(), "test.db"))
	t.Setenv("JWT_SECRET", "test")

	// jwtSecret 是包级变量，只有 LoadJWTSecret 会填充它；
	// 不调用的话签名和校验都用空 key，测试即使通过也覆盖不到真实初始化路径。
	if err := utils.LoadJWTSecret(); err != nil {
		t.Fatal(err)
	}
	if err := database.InitDB(filepath.Join(t.TempDir(), "test.db")); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.DB.Close() })

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
		t.Fatalf("Expected status code %d, but got %d, and body=%s", expected, rec.Code, rec.Body.String())
	}
}

func signupAndLogin(t *testing.T, e *gin.Engine, email string) string {
	t.Helper()
	body := fmt.Sprintf(`{"email":%q,"password":%q}`, email, testPassword)

	rec := doTest(t, e, "POST", "/signup", "", body)
	assertStatusCode(t, rec, 201)

	rec = doTest(t, e, "POST", "/login", "", body)
	assertStatusCode(t, rec, 200)

	var resp struct {
		Token string `json:"token"`
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal login response: %v, body=%s", err, rec.Body.String())
	}
	if resp.Token == "" {
		t.Fatalf("Token is empty, body=%s", rec.Body.String())
	}
	return resp.Token
}

func eventPath(id int64, suffix string) string {
	return "/events/" + strconv.FormatInt(id, 10) + suffix
}

// createTestEvent 创建一个事件并返回它的 id。
func createTestEvent(t *testing.T, e *gin.Engine, token string) int64 {
	t.Helper()

	rec := doTest(t, e, "POST", "/events", token, validEventBody)
	assertStatusCode(t, rec, 201)

	// models.Event 没有 json tag，响应里键名是 "ID"；
	// 外层再包一层 {"message":..., "event": {...}}。
	var resp struct {
		Event struct {
			ID int64 `json:"ID"`
		} `json:"event"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal event response: %v, body=%s", err, rec.Body.String())
	}
	if resp.Event.ID == 0 {
		t.Fatalf("Event id is missing, body=%s", rec.Body.String())
	}
	return resp.Event.ID
}

// myRegistrationIDs 返回 /registrations/me 里的所有事件 id。
func myRegistrationIDs(t *testing.T, e *gin.Engine, token string) []int64 {
	t.Helper()

	rec := doTest(t, e, "GET", "/registrations/me", token, "")
	assertStatusCode(t, rec, 200)

	var registrations []struct {
		ID int64 `json:"ID"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &registrations); err != nil {
		t.Fatalf("Failed to unmarshal registrations: %v, body=%s", err, rec.Body.String())
	}

	ids := make([]int64, 0, len(registrations))
	for _, registration := range registrations {
		ids = append(ids, registration.ID)
	}
	return ids
}
