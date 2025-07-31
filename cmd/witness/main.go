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
	// TODO:
	// add flags for TLS key and cert

	// Flag definitions
	port := flag.String("port", "8080", "server listen port")
	host := flag.String("host", "localhost", "host to bind to")
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

func parseTokens(tokens string) []string {
	// If no tokens, return empty slice
	if tokens == "" {
		return []string{}
	}

	return strings.Split(tokens, ",")
}
