package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	_ "github.com/fingertips18/fingertips18.github.io/backend/docs"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	v1 "github.com/fingertips18/fingertips18.github.io/backend/internal/handler/v1"
	"github.com/fingertips18/fingertips18.github.io/backend/pkg/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	httpServer *http.Server
	port       string
}

type Config struct {
	ClientURL           string
	Environment         string
	Port                string
	AuthToken           string
	EmailJSServiceID    string
	EmailJSTemplateID   string
	EmailJSPublicKey    string
	EmailJSPrivateKey   string
	GoogleMeasurementID string
	GoogleAPISecret     string
	Username            string
	Password            string
	DatabaseAPI         database.DatabaseAPI
}

type handlerConfig struct {
	paths   []string
	handler http.Handler
}

// New creates and returns a new Server instance configured with the provided Config.
// It initializes the HTTP handlers, sets up the HTTP server with the specified port and timeouts,
// and logs the environment in which the server is starting.
func New(cfg Config) *Server {
	log.Printf("Starting server with environment: %s", cfg.Environment)

	authInterceptor := middleware.NewAuthInterceptor(cfg.AuthToken)
	corsInterceptor := middleware.NewCorsInterceptor(
		middleware.CorsInterceptor{
			ClientURL: cfg.ClientURL,
			Local:     cfg.Environment == "local",
		},
	)

	handlers := createHandlers(cfg)
	mux := setupHandlers(handlers...)

	// Chain: CORS → Auth → Mux
	appChain := corsInterceptor.CorsMiddleware(
		authInterceptor.MiddlewareFunc(mux),
	)

	// Swagger UI (BasicAuth protected)
	swaggerHandler := basicAuth(cfg.Username, cfg.Password, httpSwagger.WrapHandler)

	// Top-level mux
	rootMux := http.NewServeMux()

	// Swagger path stays open (only BasicAuth)
	rootMux.Handle("/swagger/", swaggerHandler)

	// Root redirect
	rootMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger/", http.StatusFound)
			return
		}
		// everything else goes through appChain
		appChain.ServeHTTP(w, r)
	})

	return &Server{
		httpServer: &http.Server{
			Addr:              ":" + cfg.Port,
			Handler:           rootMux,
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

	projectHandler := v1.NewProjectServiceHandler(
		v1.ProjectServiceConfig{
			DatabaseAPI: cfg.DatabaseAPI,
		},
	)

	educationHandler := v1.NewEducationServiceHandler(
		v1.EducationServiceConfig{
			DatabaseAPI: cfg.DatabaseAPI,
		},
	)

	handlers := []handlerConfig{
		{
			paths:   []string{"/email", "/email/"},
			handler: emailHandler,
		},
		{
			paths:   []string{"/analytics", "/analytics/"},
			handler: analyticsHandler,
		},
		{
			paths:   []string{"/project", "/project/", "/projects", "/projects/"},
			handler: projectHandler,
		},
		{
			paths:   []string{"/education", "/education/", "/educations", "/educations/"},
			handler: educationHandler,
		},
	}

	return handlers
}

// setupHandlers registers multiple HTTP handlers to their respective paths on a new http.ServeMux.
// It accepts a variadic list of handlerConfig, where each handlerConfig contains one or more paths
// and an associated http.Handler. Each path in the handlerConfig is mapped to the provided handler.
// Returns the configured *http.ServeMux.
func setupHandlers(h ...handlerConfig) *http.ServeMux {
	mux := http.NewServeMux()

	for _, handleCfg := range h {
		for _, path := range handleCfg.paths {
			mux.Handle(path, handleCfg.handler)
		}
	}

	return mux
}

// basicAuth is a middleware that enforces HTTP Basic Authentication on incoming requests.
// It checks the provided username and password against the server's configured credentials.
// If authentication fails, it responds with a 401 Unauthorized status and a WWW-Authenticate header.
// On successful authentication, it sets an Authorization header with a Bearer token and calls the next handler.
func basicAuth(username, password string, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if !ok || u != username || p != password {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}

// Run starts the HTTP server and listens for incoming requests on the configured port.
// It logs the server startup and returns an error if the server fails to start,
// except when the error is due to the server being closed gracefully.
func (s *Server) Run() error {
	log.Printf("Starting server on http://localhost:%s", s.port)
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
