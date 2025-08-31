package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	v1 "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1"
)

type AnalyticsHandler interface {
	http.Handler
	PageView(w http.ResponseWriter, r *http.Request)
}

type AnalyticsServiceConfig struct {
	GoogleMeasurementID string
	GoogleAPISecret     string

	analyticsRepo v1.AnalyticsRepository
}

type analyticsServiceHandler struct {
	analyticsRepo v1.AnalyticsRepository
}

// NewAnalyticsServiceHandler creates a new instance of AnalyticsHandler using the provided AnalyticsServiceConfig.
// If the analytics repository is not provided in the config, it initializes a default AnalyticsRepository
// with the given Google Measurement ID and API Secret from the config.
// Returns an implementation of AnalyticsHandler.
func NewAnalyticsServiceHandler(cfg AnalyticsServiceConfig) AnalyticsHandler {
	analyticsRepo := cfg.analyticsRepo
	if analyticsRepo == nil {
		analyticsRepo = v1.NewAnalyticsRepository(
			v1.AnalyticsRepositoryConfig{
				GoogleMeasurementID: cfg.GoogleMeasurementID,
				GoogleAPISecret:     cfg.GoogleAPISecret,
			},
		)
	}

	return &analyticsServiceHandler{
		analyticsRepo: analyticsRepo,
	}
}

// ServeHTTP handles HTTP requests routed to the analytics service handler.
// It inspects the request path after trimming the "/analytics" prefix and dispatches
// the request to the appropriate handler method. Currently, it supports the "/page-view"
// endpoint by invoking the PageView handler. For any other paths, it responds with a 404 Not Found.
func (h *analyticsServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/analytics")

	switch path {
	case "/page-view":
		h.PageView(w, r)
	default:
		http.NotFound(w, r)
	}
}

// PageView handles HTTP POST requests for recording a page view analytics event.
// It expects a JSON payload in the request body containing the required fields
// PageLocation and PageTitle. If the request method is not POST, or if the JSON
// is invalid or missing required fields, it responds with an appropriate HTTP error.
// On success, it records the page view using the analytics repository and responds
// with a JSON status message.
func (h *analyticsServiceHandler) PageView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed: only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req domain.PageView
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	if req.PageLocation == "" || req.PageTitle == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := h.analyticsRepo.PageView(req.PageLocation, req.PageTitle); err != nil {
		http.Error(w, "Failed to view page: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"status": "ok"}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
