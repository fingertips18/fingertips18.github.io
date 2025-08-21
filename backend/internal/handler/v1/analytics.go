package v1

import (
	"encoding/json"
	"net/http"
	"strings"

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

type PageViewRequest struct {
	PageLocation string `json:"location"`
	PageTitle    string `json:"title"`
}

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

func (h *analyticsServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/analytics")

	switch path {
	case "/page-view":
		h.PageView(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *analyticsServiceHandler) PageView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed: only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req PageViewRequest
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]string{"status": "ok"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return
	}
}
