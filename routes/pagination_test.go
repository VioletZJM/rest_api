package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"example.com/rest_api/models"
	"github.com/gin-gonic/gin"
)

type paginatedEvents struct {
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Items []models.Event `json:"Items"`
}

// createTestEvents 连续创建 n 个事件，并返回它们的 id。
func createTestEvents(t *testing.T, e *gin.Engine, token string, n int) []int64 {
	t.Helper()

	ids := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		ids = append(ids, createTestEvent(t, e, token))
	}
	return ids
}

// fetchEvents 请求 /events 并解析分页响应。
func fetchEvents(t *testing.T, e *gin.Engine, path string) paginatedEvents {
	t.Helper()

	rec := doTest(t, e, "GET", path, "", "")
	assertStatusCode(t, rec, http.StatusOK)

	var resp paginatedEvents
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal events response: %v, body=%s", err, rec.Body.String())
	}
	return resp
}

func assertPagination(t *testing.T, got paginatedEvents, wantTotal, wantPage, wantLimit, wantItems int) {
	t.Helper()

	if got.Total != wantTotal {
		t.Errorf("Expected total %d, but got %d", wantTotal, got.Total)
	}
	if got.Page != wantPage {
		t.Errorf("Expected page %d, but got %d", wantPage, got.Page)
	}
	if got.Limit != wantLimit {
		t.Errorf("Expected limit %d, but got %d", wantLimit, got.Limit)
	}
	if len(got.Items) != wantItems {
		t.Errorf("Expected %d items, but got %d", wantItems, len(got.Items))
	}
}

func TestGetEventsPagination(t *testing.T) {
	e := newTestServer(t)
	token := signupAndLogin(t, e, "pagination@example.com")
	eventIDs := createTestEvents(t, e, token, 25)

	created := make(map[int64]bool, len(eventIDs))
	for _, id := range eventIDs {
		created[id] = true
	}

	t.Run("first page returns 20 of 25 events", func(t *testing.T) {
		got := fetchEvents(t, e, "/events?page=1&limit=20")
		assertPagination(t, got, 25, 1, 20, 20)
	})

	t.Run("second page returns the remaining 5 events", func(t *testing.T) {
		got := fetchEvents(t, e, "/events?page=2&limit=20")
		assertPagination(t, got, 25, 2, 20, 5)
	})

	t.Run("pages do not overlap and cover all events", func(t *testing.T) {
		first := fetchEvents(t, e, "/events?page=1&limit=20")
		second := fetchEvents(t, e, "/events?page=2&limit=20")

		seen := make(map[int64]bool, 25)
		for _, event := range append(first.Items, second.Items...) {
			if !created[event.ID] {
				t.Errorf("Unexpected event id %d", event.ID)
			}
			if seen[event.ID] {
				t.Errorf("Event id %d appears on both pages", event.ID)
			}
			seen[event.ID] = true
		}
		if len(seen) != 25 {
			t.Errorf("Expected 25 distinct events across both pages, but got %d", len(seen))
		}
	})

	t.Run("defaults to page 1 with limit 20", func(t *testing.T) {
		got := fetchEvents(t, e, "/events")
		assertPagination(t, got, 25, 1, 20, 20)
	})

	t.Run("limit above 100 is clamped to 100", func(t *testing.T) {
		for _, limit := range []string{"101", "1000"} {
			t.Run(limit, func(t *testing.T) {
				got := fetchEvents(t, e, "/events?limit="+limit)
				assertPagination(t, got, 25, 1, 100, 25)
			})
		}
	})

	t.Run("limit below 1 falls back to 20", func(t *testing.T) {
		for _, limit := range []string{"0", "-1", "abc"} {
			t.Run(limit, func(t *testing.T) {
				got := fetchEvents(t, e, "/events?limit="+limit)
				assertPagination(t, got, 25, 1, 20, 20)
			})
		}
	})

	t.Run("non-integer page falls back to page 1", func(t *testing.T) {
		for _, page := range []string{"abc", "1.5", ""} {
			t.Run(fmt.Sprintf("page=%q", page), func(t *testing.T) {
				got := fetchEvents(t, e, "/events?page="+page)
				assertPagination(t, got, 25, 1, 20, 20)
			})
		}
	})

	t.Run("page beyond last page returns empty list instead of 404", func(t *testing.T) {
		got := fetchEvents(t, e, "/events?page=99")
		assertPagination(t, got, 25, 99, 20, 0)
	})
}
