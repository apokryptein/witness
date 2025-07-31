package middleware

import (
	"net/http"
	"slices"
	"strings"
)

// RequireAuth is a middleware http.HandlerFunc that requires authorization
// for an existing http.Handler route when applied
// Expects a standard bearertoken format: "Authorization: Bearer TOKEN"
func RequireAuth(validTokens []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If not valid tokens exist, we are defaulting to no auth
		if len(validTokens) == 0 {
			next(w, r)
			return
		}

		// Retrieve and verify auth header
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Retrieve token from auth header & verify
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if slices.Contains(validTokens, token) {
			next(w, r)
			return
		}

		// If token doesn't match a valid token, return Unauthorized
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
