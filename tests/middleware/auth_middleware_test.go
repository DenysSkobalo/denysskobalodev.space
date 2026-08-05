package handlers_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DenysSkobalo/denysskobalodev.space/internal/handlers"
	"golang.org/x/crypto/pbkdf2"
)

type contextKey string

type MockKV struct {
	data map[string]string
}

func (m *MockKV) Put(k, v string, ttl int) error { m.data[k] = v; return nil }
func (m *MockKV) Delete(k string) error         { delete(m.data, k); return nil }
func (m *MockKV) Get(k string) (string, error)  { return m.data[k], nil }

func TestLoginWorkflow(t *testing.T) {
	salt := "denysskobalo_unique_salt"
	password := "SecretPassword123!"

	rawHash := pbkdf2.Key([]byte(password), []byte(salt), 100000, 32, sha256.New)
	expectedHash := hex.EncodeToString(rawHash)

	kv := &MockKV{data: make(map[string]string)}
	handler := &handlers.AuthHandler{
		AdminHash: expectedHash,
		Salt:      salt,
		KVStore:   kv,
	}

	t.Run("Invalid_Credentials", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"username": "admin", "password": "WrongPassword"})
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		
		ctx := context.WithValue(req.Context(), contextKey("ADMIN_HASH"), expectedHash)
		ctx = context.WithValue(ctx, contextKey("SALT"), salt)
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()
		handler.Login(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401 Unauthorized, got %d", rec.Code)
		}
	})

	t.Run("Successful_Login_and_Cookie_Placement", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"username": "admin", "password": password})
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))

		ctx := context.WithValue(req.Context(), contextKey("ADMIN_HASH"), expectedHash)
		ctx = context.WithValue(ctx, contextKey("SALT"), salt)
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()
		handler.Login(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected 200 OK, got %d", rec.Code)
		}

		cookies := rec.Result().Cookies()
		found := false
		for _, c := range cookies {
			if c.Name == "admin_session" && c.Value != "" {
				found = true
				normalizedDomain := strings.TrimPrefix(c.Domain, ".")
				expectedDomain := strings.TrimPrefix(".denysskobalodev.space", ".")

				if normalizedDomain != expectedDomain {
					t.Errorf("Expected domain %s, got %s", expectedDomain, c.Domain)
				}
			}
		}
		if !found {
			t.Error("admin_session cookie was not properly set in response headers")
		}
	})
}
