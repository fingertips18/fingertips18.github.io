package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	v1 "github.com/Fingertips18/fingertips18.github.io/backend/internal/handler/v1"
	"github.com/Fingertips18/fingertips18.github.io/backend/pkg/middleware"
)

type Server struct {
	httpServer *http.Server
	port       string
}

type Config struct {
	Environment         string
	Port                string
	AuthToken           string
	EmailJSServiceID    string
	EmailJSTemplateID   string
	EmailJSPublicKey    string
	EmailJSPrivateKey   string
	GoogleMeasurementID string
	GoogleAPISecret     string
}

type handlerConfig struct {
	path    string
	handler http.Handler
}

// New creates and returns a new Server instance configured with the provided Config.
// It initializes the HTTP handlers, sets up the HTTP server with the specified port and timeouts,
// and logs the environment in which the server is starting.
func New(cfg Config) *Server {
	log.Printf("Starting server with environment: %s", cfg.Environment)

	authInterceptor := middleware.NewAuthInterceptor(cfg.AuthToken)

	handlers := createHandlers(cfg)
	mux := setupHandlers(handlers...)

	muxWithAuth := http.NewServeMux()

	muxWithAuth.Handle("/", authInterceptor.MiddlewareFunc(mux))

	return &Server{
		httpServer: &http.Server{
			Addr:              ":" + cfg.Port,
			Handler:           muxWithAuth,
			ReadHeaderTimeout: 30 * time.Second,
		},
		port: cfg.Port,
	}
}

// createHandlers initializes and returns a slice of handlerConfig structs,
// each representing an HTTP handler for the server. It configures the handlers
// using the provided Config, such as setting up the email service handler with
// the necessary credentials and service IDs.
func createHandlers(cfg Config) []handlerConfig {
	emailHandler := v1.NewEmailServiceHandler(
		v1.EmailServiceConfig{
			ServiceID:   cfg.EmailJSServiceID,
			TemplateID:  cfg.EmailJSTemplateID,
			UserID:      cfg.EmailJSPublicKey,
			AccessToken: cfg.EmailJSPrivateKey,
		},
	)

	analyticsHandler := v1.NewAnalyticsServiceHandler(
		v1.AnalyticsServiceConfig{
			GoogleMeasurementID: cfg.GoogleMeasurementID,
			GoogleAPISecret:     cfg.GoogleAPISecret,
		},
	)

	handlers := []handlerConfig{
		{
			path:    "/email/",
			handler: emailHandler,
		},
		{
			path:    "/analytics/",
			handler: analyticsHandler,
		},
	}

	return handlers
}

// setupHandlers creates and returns a new http.ServeMux with the provided handler configurations registered.
// Each handlerConfig in the variadic parameter 'h' specifies a path and its corresponding handler to be registered
// with the ServeMux. This function is useful for setting up multiple route handlers in a concise manner.
//
// Parameters:
//
//	h ...handlerConfig - A variadic list of handlerConfig structs, each containing a path and an http.Handler.
//
// Returns:
//
//	*http.ServeMux - A pointer to the configured http.ServeMux instance.
func setupHandlers(h ...handlerConfig) *http.ServeMux {
	mux := http.NewServeMux()

	for _, handleCfg := range h {
		mux.Handle(handleCfg.path, handleCfg.handler)
	}

	return mux
}

// Run starts the HTTP server and listens for incoming requests on the configured port.
// It logs the server startup and returns an error if the server fails to start,
// except when the error is due to the server being closed gracefully.
func (s *Server) Run() error {
	log.Printf("Starting server on port %s", s.port)
	// Listen on port 8080
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Could not listen on %s: %v", s.port, err)
	}

	return nil
}

// Shutdown gracefully shuts down the HTTP server associated with the Server instance.
// It attempts to shut down the server using the provided context, allowing for any
// in-flight requests to complete before termination. If an error occurs during shutdown,
// it logs the error and returns it.
func (s *Server) Shutdown(ctx context.Context) error {
	// Shutdown the HTTP server
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Error shutting down HTTP server: %v", err)
	}

	return err
}
