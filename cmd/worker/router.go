//go:build js && wasm

package main

import (
	"net/http"
	"github.com/DenysSkobalo/denysskobalodev.space/internal/handlers"
	"github.com/DenysSkobalo/denysskobalodev.space/internal/middleware"
)

type KVPutOp struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

type KVDelOp struct {
	Key string `json:"key"`
}

var (
	PendingKVPut      *KVPutOp
	PendingKVDelete   *KVDelOp
	PrefetchedSession string
)

type BridgeKV struct{}

func (b *BridgeKV) Put(key, val string, ttl int) error {
	PendingKVPut = &KVPutOp{Key: key, Value: val, TTL: ttl}
	return nil
}

func (b *BridgeKV) Delete(key string) error {
	PendingKVDelete = &KVDelOp{Key: key}
	return nil
}

func (b *BridgeKV) Get(key string) (string, error) {
	return PrefetchedSession, nil 
}

func setupRouter() http.Handler {
	mux := http.NewServeMux()

	authHandler := &handlers.AuthHandler{
		KVStore: &BridgeKV{},
	}
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
