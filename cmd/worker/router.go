//go:build js && wasm

package main

import (
	"net/http"

	"github.com/DenysSkobalo/denysskobalodev.space/internal/handlers"
	"github.com/DenysSkobalo/denysskobalodev.space/internal/middleware"
)

func setupRouter() http.Handler {
	mux := http.NewServeMux()

	authHandler := &handlers.AuthHandler{}
	projectsHandler := &handlers.ProjectsHandler{}

	mux.HandleFunc("GET /api/telemetry", handlers.TelemetryHandler)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("GET /api/projects", projectsHandler.List)

	adminAuth := middleware.RequireAdmin(authHandler)

	mux.Handle("GET /api/auth/me", adminAuth(http.HandlerFunc(authHandler.Me)))
	mux.Handle("POST /api/auth/logout", adminAuth(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("POST /api/projects", adminAuth(http.HandlerFunc(projectsHandler.Create)))

	return mux
}
