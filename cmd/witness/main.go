package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/apokryptein/witness/internal/handler"
	"github.com/apokryptein/witness/internal/middleware"
	"github.com/apokryptein/witness/internal/server"
)

func main() {
	// Flag definitions
	port := flag.String("port", "8443", "server listen port")
	host := flag.String("host", "localhost", "host to bind to")
	tlsCert := flag.String("tls-cert", "", "path to TLS certificate file")
	tlsKey := flag.String("tls-key", "", "path to TLS key file")
	tokens := flag.String("tokens", "", "bearer tokens (comma-separated, none = no auth)")
	flag.Parse()

	// Setup context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse bearer tokens
	bearerTokens := parseTokens(*tokens)

	// Create server config
	s := server.Server{
		ListenAddr: fmt.Sprintf("%s:%s", *host, *port),
		TLSCert:    *tlsCert,
		TLSKey:     *tlsKey,
		Handler:    setupRoutes(bearerTokens),
	}

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Serve
	go func() {
		if err := s.Run(ctx); err != nil {
			log.Printf("[ERR] server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("[INFO] Server is shutting down...")
	cancel()
}

// setupRoutes instantiates a new http.ServeMux and instantiates our routes with the desired
// middleware
func setupRoutes(tokens []string) *http.ServeMux {
	// Instantiate new mux
	mux := http.NewServeMux()

	// Add routes
	mux.HandleFunc("/echo",
		middleware.WithCORS([]string{"GET"})(
			middleware.WithLogging(
				middleware.RequireAuth(tokens, handler.EchoHandler))))

	mux.HandleFunc("/ip",
		middleware.WithCORS([]string{"GET"})(
			middleware.WithLogging(
				middleware.RequireAuth(tokens, handler.IPHandler))))

	mux.HandleFunc("/health",
		middleware.WithCORS([]string{"GET"})(
			middleware.WithLogging(
				middleware.RequireAuth(tokens, handler.HealthHandler))))

	mux.HandleFunc("/whoami",
		middleware.WithCORS([]string{"GET"})(
			middleware.WithLogging(
				middleware.RequireAuth(tokens, handler.WhoHandler))))

	mux.HandleFunc("/headers",
		middleware.WithCORS([]string{"GET"})(
			middleware.WithLogging(
				middleware.RequireAuth(tokens, handler.HeaderHandler))))

	mux.HandleFunc("/",
		middleware.WithCORS([]string{"GET"})(
			middleware.WithLogging(
				middleware.RequireAuth(tokens, handler.NotFoundHandler))))

	return mux
}

// parseTokens parses a comma-separated list of bearer tokens and returns them
// in a string slice
func parseTokens(tokens string) []string {
	// If no tokens, return empty slice
	if tokens == "" {
		return []string{}
	}

	return strings.Split(tokens, ",")
}
