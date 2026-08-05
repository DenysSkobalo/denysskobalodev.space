package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DenysSkobalo/denysskobalodev.space/internal/handlers"
	"github.com/DenysSkobalo/denysskobalodev.space/internal/middleware"
)

// MockKVClient реалізує KVClient для локальної розробки / тестування
type MockKVClient struct {
	store map[string]string
}

func (m *MockKVClient) Put(key string, val string, ttl int) error {
	m.store[key] = val
	return nil
}
func (m *MockKVClient) Delete(key string) error {
	delete(m.store, key)
	return nil
}
func (m *MockKVClient) Get(key string) (string, error) {
	val, ok := m.store[key]
	if !ok {
		return "", nil
	}
	return val, nil
}

func main() {
	allowedOrigin := "https://denysskobalodev.space"
	adminHash := os.Getenv("ADMIN_PASSWORD_HASH")
	salt := "denysskobalo_unique_salt"

	// Локальна емуляція KV (у продакшені замінюється на Cloudflare KV Binding)
	mockKV := &MockKVClient{store: make(map[string]string)}

	mux := http.NewServeMux()

	// 1. Telemetry
	mux.HandleFunc("GET /api/telemetry", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"experienceYears":"5+","dailyEvents":"5M+","avgLatencyMs":"<50ms","activeSessions":"3,000+","systemStatus":"OPERATIONAL"}`))
	})

	// 2. Auth Handlers
	authHandler := &handlers.AuthHandler{
		AdminHash: adminHash,
		Salt:      salt,
		KVStore:   mockKV,
	}

	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)
	mux.HandleFunc("GET /api/auth/me", authHandler.Me)

	// 3. Protected Projects API
	projHandler := &handlers.ProjectsHandler{}
	mux.HandleFunc("GET /api/projects", projHandler.List)

	// Auth Guard Middleware
	adminGuard := middleware.RequireAdmin(authHandler)
	mux.Handle("POST /api/projects", adminGuard(http.HandlerFunc(projHandler.Create)))

	// Global Middleware Pipeline
	handler := middleware.CORS(allowedOrigin)(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[Production Go Engine] Server initialized on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server launch failed: %v", err)
	}
}
