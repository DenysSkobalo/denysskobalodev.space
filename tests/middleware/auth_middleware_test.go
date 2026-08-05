package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DenysSkobalo/denysskobalodev.space/internal/middleware"
)

type MockValidator struct {
	validToken string
}

func (m *MockValidator) ValidateSession(token string) bool {
	return token == m.validToken
}

func TestRequireAdminMiddleware(t *testing.T) {
	validator := &MockValidator{validToken: "valid_secret_session_token"}
	guard := middleware.RequireAdmin(validator)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ACCESS_GRANTED"))
	})

	handlerToTest := guard(nextHandler)

	t.Run("Should Return 401 Unauthorized When Cookie Missing", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/projects", nil)
		rec := httptest.NewRecorder()

		handlerToTest.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", rec.Code)
		}
	})

	t.Run("Should Pass When Valid Cookie Provided", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/projects", nil)
		req.AddCookie(&http.Cookie{Name: "admin_session", Value: "valid_secret_session_token"})
		rec := httptest.NewRecorder()

		handlerToTest.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
	})
}
