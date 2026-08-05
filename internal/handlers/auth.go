package handlers

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

type AuthHandler struct {
	AdminHash string
	Salt      string
	KVStore   KVClient
}

type KVClient interface {
	Put(key string, val string, ttl int) error
	Delete(key string) error
	Get(key string) (string, error)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Username and password are required")
		return
	}

	// PBKDF2 Hash Calculation
	computedHash := pbkdf2.Key([]byte(req.Password), []byte(h.Salt), 100000, 32, sha256.New)
	computedHex := hex.EncodeToString(computedHash)

	// Constant Time Comparison (Timing Attack Protection)
	if req.Username != "admin" || subtle.ConstantTimeCompare([]byte(computedHex), []byte(h.AdminHash)) != 1 {
		respondError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid username or password")
		return
	}

	// Generate Crypto Session Token
	sessionID := generateCryptoToken()
	_ = h.KVStore.Put("session:"+sessionID, `{"user":"admin","role":"administrator"}`, 86400)

	// 4. Set HttpOnly Cookie for Unified Domain
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    sessionID,
		Path:     "/",
		Domain:   ".denysskobalodev.space",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Successfully authenticated",
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("admin_session"); err == nil {
		_ = h.KVStore.Delete("session:" + cookie.Value)
	}

	// Invalidate Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/",
		Domain:   ".denysskobalodev.space",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("admin_session")
	if err != nil || cookie.Value == "" {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Session expired or invalid")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"user": map[string]string{
			"username": "admin",
			"role":     "administrator",
		},
	})
}

func (h *AuthHandler) ValidateSession(token string) bool {
	if token == "" || h.KVStore == nil {
		return false
	}
	val, err := h.KVStore.Get("session:" + token)
	if err != nil || val == "" {
		return false
	}
	return true
}
