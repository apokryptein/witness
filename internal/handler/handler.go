// Package handler contains all types, functions, and methods for Witness
// HTTP handlers
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ========================= //
// === STRUCT DEFINITIONS == //
// ========================= //

// TLSData represents pertinent TLS data
type TLSData struct {
	Version     uint16
	CipherSuite uint16
}

// Whoami represents selected data for the whoami endpoint
type Whoami struct {
	IP      string
	TLS     TLSData
	Headers map[string]string
}

// ========================== //
// === HANDLER DEFINITIONS == //
// ========================== //

// EchoHandler is an http.Handler containing logic for simple echo
// server functionality
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
	}

	// Set content-length header for echo reply
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	w.WriteHeader(http.StatusOK)

	// Echo body back
	w.Write(body)
}

// IPHandler is an http.Handler that contains logic for an endpoint to return the client's IP address
func IPHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch IP
	ip := fetchIP(r)

	// Set and write applicable headers
	w.Header().Set("Content-Type", "text/plain; chartset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(ip)))
	w.WriteHeader(http.StatusOK)

	// Return IP address
	w.Write([]byte(ip))
}

// HealthHandler is an http.Handler that returns a status and timestamp
// if the server is running and functional
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Set and write headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Build response
	response := map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}

// HeaderHandler is an http.Handler that returns all headers for a given request
// to the client
func HeaderHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch headers
	headers := fetchHeaders(r)

	// Set and write response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(headers); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}

// WhoHandler is an http.Handler that returns various data back to the requester:
// IP address, TLS info, and request headers
func WhoHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch relevant data
	ip := fetchIP(r)
	tlsInfo := fetchTLS(r)
	headers := fetchHeaders(r)

	// Set and write response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Build whoami response
	whoami := Whoami{
		IP:      ip,
		TLS:     tlsInfo,
		Headers: headers,
	}

	if err := json.NewEncoder(w).Encode(whoami); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}

// NotFoundHandler is an http.Handler that is called when an invalid server route
// is called
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotFound)
}

// =================================== //
// === HELPER FUNCTION DEFINITIONS === //
// =================================== //

// fetchIP retrieves the X-FORWARDED-FOR header to check
// for public IP address. If this value is nil, it returns
// the RemoteAddr instead
func fetchIP(req *http.Request) string {
	// Look for X-Forwarded-For header
	// X-Forwarded-For: <client>, <proxy>
	// X-Forwarded-For: <client>, <proxy>, …, <proxyN>
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP (original client)
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}

		// If only one value (no proxies), return XFF
		return strings.TrimSpace(xff)
	}

	// Check for X-Real-IP (nginx)
	if xri := req.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Otherwise, return req.RemoteAddr
	return req.RemoteAddr
}

// fetchHeaders retrieves and retruns all headers from the provided
// HTTP request
func fetchHeaders(req *http.Request) map[string]string {
	headers := make(map[string]string)

	// Iterate over all headers and store in headers map
	for header, values := range req.Header {
		// If more than one value join with ','
		headers[header] = strings.Join(values, ", ")
	}

	return headers
}

// fetchTLS retrieves relevant TLS data (version and ciphersuite)
// if HTTPS/TLS is used
func fetchTLS(req *http.Request) TLSData {
	if req.TLS == nil {
		return TLSData{}
	}

	tlsInfo := req.TLS

	return TLSData{
		Version:     tlsInfo.Version,
		CipherSuite: tlsInfo.CipherSuite,
	}
}
