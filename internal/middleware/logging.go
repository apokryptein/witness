// Package middleware contains all middleware server functionality
package middleware

import (
	"log"
	"net/http"
)

// WithLogging is a middleware http.HandlerFunc that adds logging to
// an existing http.Handler
func WithLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}
