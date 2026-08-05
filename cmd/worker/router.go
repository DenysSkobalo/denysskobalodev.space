package main

import (
	"net/http"
	"sync"

	"github.com/DenysSkobalo/denysskobalodev.space/internal/handlers"
	"github.com/DenysSkobalo/denysskobalodev.space/internal/middleware"
)

type MemoryKV struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMemoryKV() *MemoryKV {
	return &MemoryKV{data: make(map[string]string)}
}

func (m *MemoryKV) Put(key string, val string, ttl int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
	return nil
}

func (m *MemoryKV) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

func (m *MemoryKV) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key], nil
}

func setupRouter() http.Handler {
	allowedOrigin := "https://denysskobalodev.space"

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/telemetry", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"experienceYears":"5+","dailyEvents":"5M+","avgLatencyMs":"<50ms","activeSessions":"3,000+","systemStatus":"OPERATIONAL"}`))
	})

	memKV := NewMemoryKV()
	authHandler := &handlers.AuthHandler{
		AdminHash: "32439c09191d8424a2ef693f185ef30e58f00a5d4d3e8e202166567302482e9b", // PBKDF2 hash
		Salt:      "denysskobalo_unique_salt",
		KVStore:   memKV,
	}

	projHandler := &handlers.ProjectsHandler{
		DB: nil,
	}

	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)
	mux.HandleFunc("GET /api/auth/me", authHandler.Me)

	mux.HandleFunc("GET /api/projects", projHandler.List)
	mux.HandleFunc("POST /api/projects", projHandler.Create)

	return middleware.CORS(allowedOrigin)(mux)
}
