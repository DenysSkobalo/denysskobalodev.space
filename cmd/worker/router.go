package main

import (
	"net/http"

	"github.com/DenysSkobalo/denysskobalodev.space/internal/handlers"
	"github.com/DenysSkobalo/denysskobalodev.space/internal/middleware"
)

func setupRouter() http.Handler {
	allowedOrigin := "https://denysskobalodev.space"

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/telemetry", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"experienceYears":"5+","dailyEvents":"5M+","avgLatencyMs":"<50ms","activeSessions":"3,000+","systemStatus":"OPERATIONAL"}`))
	})

	authHandler := &handlers.AuthHandler{}
	projHandler := &handlers.ProjectsHandler{}

	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)
	mux.HandleFunc("GET /api/auth/me", authHandler.Me)
	mux.HandleFunc("GET /api/projects", projHandler.List)

	return middleware.CORS(allowedOrigin)(mux)
}
