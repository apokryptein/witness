// Package server is responsible for the configuration and running
// of the HTTP/S server
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Server represents an HTTP/S server
type Server struct {
	ListenAddr string       // IP:PORT
	TLSCert    string       // Path to TLS certificate
	TLSKey     string       // Path to TLS key
	Handler    http.Handler // *http.ServeMux containing our routes
}

// Run runs an HTTP/S server for a given server configuration struct
func (s *Server) Run(ctx context.Context) error {
	// Instantiate http.Server w/ relevant values
	server := &http.Server{
		Addr:    s.ListenAddr,
		Handler: s.Handler,
	}

	// DEBUG
	fmt.Printf("[INFO] Server listening on: %s\n", s.ListenAddr)

	// Create error channel
	errChan := make(chan error, 1)

	// Run server in go func -> non-blocking
	go func() {
		var err error
		if s.TLSCert != "" && s.TLSKey != "" {
			err = server.ListenAndServeTLS(s.TLSCert, s.TLSKey)
		} else {
			err = server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Listen for an error from the server or from a sigterm or interrupt
	select {
	case err := <-errChan:
		return fmt.Errorf("server failed: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}
