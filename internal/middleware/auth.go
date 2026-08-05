package middleware

import "net/http"

type SessionValidator interface {
	ValidateSession(token string) bool
}

func RequireAdmin(validator SessionValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("admin_session")
			if err != nil || cookie.Value == "" || !validator.ValidateSession(cookie.Value) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":{"code":"UNAUTHORIZED","message":"Invalid or expired session credentials"}}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}


