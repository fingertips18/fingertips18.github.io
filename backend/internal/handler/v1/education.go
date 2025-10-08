package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	v1 "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1"
)

type EducationHandler interface {
	http.Handler
	Create(w http.ResponseWriter, r *http.Request)
}

type EducationServiceConfig struct {
	DatabaseAPI database.DatabaseAPI

	educationRepo v1.EducationRepository
}

type educationServiceHandler struct {
	educationRepo v1.EducationRepository
}

// NewEducationServiceHandler creates and returns an EducationHandler configured using the provided
// EducationServiceConfig. If cfg.educationRepo is nil, a default repository is constructed via
// v1.NewEducationRepository using cfg.DatabaseAPI and the "Education" table. The returned handler
// wraps the chosen repository and is ready to serve education-related operations.
func NewEducationServiceHandler(cfg EducationServiceConfig) EducationHandler {
	educationRepo := cfg.educationRepo
	if educationRepo == nil {
		educationRepo = v1.NewEducationRepository(
			v1.EducationRepositoryConfig{
				DatabaseAPI:    cfg.DatabaseAPI,
				EducationTable: "Education",
			},
		)
	}

	return &educationServiceHandler{
		educationRepo: educationRepo,
	}
}

func (h *educationServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Normalize path by trimming trailing slash
	path := strings.TrimSuffix(r.URL.Path, "/")

	switch path {
	// POST / PUT /education
	case "/education":
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	// Unknown route
	default:
		http.NotFound(w, r)
		return
	}
}

// Create handles HTTP POST requests to create a new education record.
// It only supports the POST method; other methods receive a 405 Method Not Allowed.
// The handler decodes a JSON request body into a domain.CreateEducation and
// defers closing the request body. It maps the decoded payload to a domain.Education
// (populating MainSchool, SchoolPeriods, Projects and Level) and calls
// h.educationRepo.Create with the request context. On success it returns a JSON
// body containing the newly created ID (domain.EducationIDResponse) with
// Content-Type "application/json" and HTTP status 201 Created. If JSON decoding
// fails the handler responds with 400 Bad Request; if creation or response
// encoding fails it responds with 500 Internal Server Error.
//
// @Security ApiKeyAuth
// @Summary Create a education
// @Description Creates a new education from the provided JSON payload. Returns the created education with an assigned ID.
// @Tags education
// @Accept json
// @Produce json
// @Param education body domain.CreateEducation true "Education payload"
// @Success 201 {string} domain.EducationIDResponse "Education ID"
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /education [post]
func (h *educationServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed: only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var createReq domain.CreateEducation
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	// Map to Education and validate BEFORE calling repository
	education := &domain.Education{
		MainSchool:    createReq.MainSchool,
		SchoolPeriods: createReq.SchoolPeriods,
		Projects:      createReq.Projects,
		Level:         createReq.Level,
	}

	// Add validation here
	if err := education.ValidatePayload(); err != nil {
		http.Error(w, "Invalid education payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.educationRepo.Create(r.Context(), education)

	if err != nil {
		http.Error(w, "Failed to create education: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := domain.EducationIDResponse{ID: id}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}
