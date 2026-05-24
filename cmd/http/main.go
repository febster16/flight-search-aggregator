package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"flight-search-aggregator/config"

	"github.com/gorilla/mux"
)

func main() {
	ctx := context.Background()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Server panicked: %v\n%s", r, string(debug.Stack()))
			os.Exit(1)
		}
	}()

	environmentFlag := os.Getenv("ENVIRONMENT")
	cfg := config.LoadConfig(ctx, environmentFlag)

	httpServer, err := NewHTTPServer(ctx, cfg)
	if err != nil {
		log.Printf("Error while initializing routes: %v", err)
		os.Exit(1)
	}

	log.Printf("Starting HTTP server [%s] at port: %d", cfg.HTTP.Name, cfg.HTTP.Port)

	err = httpServer.Start(ctx)
	if err != nil {
		log.Printf("Error on starting up HTTP Server: %v", err)
		os.Exit(1)
	}
}

// Server wraps the HTTP server.
type Server struct {
	httpServer *http.Server
	router     *mux.Router
	config     *config.HTTPConfig
}

// NewHTTPServer creates a new HTTP server with the given config.
func NewHTTPServer(ctx context.Context, conf *config.Config) (*Server, error) {
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": conf.HTTP.Name,
			"version": conf.HTTP.Version,
		})
	}).Methods(http.MethodGet)

	// TODO: Register controllers here

	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", conf.HTTP.Port),
			Handler: router,
		},
		router: router,
		config: &conf.HTTP,
	}, nil
}

// Start begins listening on the configured port.
func (s *Server) Start(ctx context.Context) error {
	log.Printf("HTTP server [%s] listening on port %d", s.config.Name, s.config.Port)
	return s.httpServer.ListenAndServe()
}
