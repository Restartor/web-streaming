package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()

	t.Run("films endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/films", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rr.Code)
		}
	})

	t.Run("login only allows post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/login", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("expected status 405, got %d", rr.Code)
		}
	})

	t.Run("login requires valid JSON body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("not-json"))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rr.Code)
		}
	})
}
