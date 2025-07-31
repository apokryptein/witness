package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/apokryptein/witness/internal/handler"
)

func main() {
	// Flag definitions
	port := flag.String("port", "8080", "server listen port")
	host := flag.String("host", "localhost", "host to bind to")
	flag.Parse()

	// Instantiate new mux
	mux := http.NewServeMux()

	// Add routes
	mux.HandleFunc("/echo", handler.EchoHandler)
	mux.HandleFunc("/ip", handler.IPHandler)
	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/whoami", handler.WhoHandler)
	mux.HandleFunc("/headers", handler.HeaderHandler)
	mux.HandleFunc("/", handler.NotFoundHandler)

	// LOG
	fmt.Printf("[INFO] Server listening on: %s:%s\n", *host, *port)

	// Serve
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", *host, *port), mux)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERR] server failed: %v\n", err)
		os.Exit(1)
	}
}
