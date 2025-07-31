// Package server is responsible for the configuration and running
// of the HTTP/S server
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	ListenAddr string
	TLSCert    string
	TLSKey     string
	Handler    http.Handler
}

func (s *Server) Run(ctx context.Context) error {
	// Instantiate http.Server w/ relevant values
	server := &http.Server{
		Addr:    s.ListenAddr,
		Handler: s.Handler,
	}

	// DEBUG
	fmt.Printf("[INFO] Server listening on: %s\n", s.ListenAddr)

	errChan := make(chan error, 1)

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

	select {
	case err := <-errChan:
		return fmt.Errorf("server failed: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}
