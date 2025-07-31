package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/apokryptein/witness/internal/handler"
	"github.com/apokryptein/witness/internal/middleware"
)

func main() {
	// Flag definitions
	port := flag.String("port", "8080", "server listen port")
	host := flag.String("host", "localhost", "host to bind to")
	tokens := flag.String("tokens", "", "bearer tokens (comma-separated, none = no auth)")
	flag.Parse()

	// Parse bearer tokens
	bearerTokens := parseTokens(*tokens)

	// Instantiate new mux
	mux := http.NewServeMux()

	// Add routes
	mux.HandleFunc("/echo",
		middleware.WithCORS([]string{"GET"})(
			middleware.WithLogging(
				middleware.RequireAuth(bearerTokens, handler.EchoHandler))))

	mux.HandleFunc("/ip",
		middleware.WithCORS([]string{"GET"})(
			middleware.WithLogging(
				middleware.RequireAuth(bearerTokens, handler.IPHandler))))

	mux.HandleFunc("/health",
		middleware.WithCORS([]string{"GET"})(
			middleware.WithLogging(
				middleware.RequireAuth(bearerTokens, handler.HealthHandler))))

	mux.HandleFunc("/whoami",
		middleware.WithCORS([]string{"GET"})(
			middleware.WithLogging(
				middleware.RequireAuth(bearerTokens, handler.WhoHandler))))

	mux.HandleFunc("/headers",
		middleware.WithCORS([]string{"GET"})(
			middleware.WithLogging(
				middleware.RequireAuth(bearerTokens, handler.HeaderHandler))))

	mux.HandleFunc("/",
		middleware.WithCORS([]string{"GET"})(
			middleware.WithLogging(
				middleware.RequireAuth(bearerTokens, handler.NotFoundHandler))))

	// LOG
	fmt.Printf("[INFO] Server listening on: %s:%s\n", *host, *port)

	// Serve
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", *host, *port), mux)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERR] server failed: %v\n", err)
		os.Exit(1)
	}
}

func parseTokens(tokens string) []string {
	// If no tokens, return empty slice
	if tokens == "" {
		return []string{}
	}

	return strings.Split(tokens, ",")
}
