package middleware

import (
	"net/http"
	"strings"
)

// WithCORS adds CORS headers
func WithCORS(allowedMethods []string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Allow all origins for development/testing
			w.Header().Set("Access-Control-Allow-Origin", "*")

			// Allow supported methods
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))

			// Allow common headers
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle pre-flight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next(w, r)
		}
	}
}
