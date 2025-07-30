// Package handler contains all types, functions, and methods for Witness
// HTTP handlers
package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// RequestHandler implements the http.Handler interface and acts
// as a very basic
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	// Logging
	log.Printf("%s %s", r.Method, r.URL.Path)

	// Switch on path
	switch r.URL.Path {
	case "/echo":
		// Read body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
		}

		// Echo body back
		fmt.Fprintf(w, "%s", body)

	case "/ip":
		// Look for X-FORWARDED-FOR header
		ip := r.Header.Get("X-FORWARDED-FOR")

		// If header not present, set IP to RemoteAddr
		if ip == "" {
			ip = r.RemoteAddr
		}

		// Return IP address
		fmt.Fprintf(w, "%s", ip)

	default:
		http.NotFound(w, r)
	}
}
