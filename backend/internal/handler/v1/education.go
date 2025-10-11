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
	Get(w http.ResponseWriter, r *http.Request, id string)
	Update(w http.ResponseWriter, r *http.Request)
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

	switch {
	// POST / PUT /education
	case path == "/education":
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodPut:
			h.Update(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		// GET / DELETE /education/{id}
	case strings.HasPrefix(path, "/education/"):
		id := strings.TrimPrefix(path, "/education/")

		if id == "" {
			http.Error(w, "Education ID is required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.Get(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return

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

// Get handles HTTP GET requests for an education resource identified by id.
// It enforces the GET method and returns StatusMethodNotAllowed for other HTTP methods.
// The handler uses the request context to retrieve the education entity from the repository.
// If the repository returns an error, Get responds with StatusInternalServerError.
// If the requested education resource is not found, Get responds with StatusNotFound.
// On success it encodes the education entity as JSON, sets Content-Type to application/json,
// and writes the payload with StatusOK. Any encoding or write error results in an internal server error response.
//
// @Security ApiKeyAuth
// @Summary Get a education by ID
// @Description Retrieves the details of a specific education using its unique ID.
// @Tags education
// @Accept json
// @Produce json
// @Param id path string true "Education ID"
// @Success 200 {object} domain.Education
// @Failure 400 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /education/{id} [get]
func (h *educationServiceHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed: only GET is supported", http.StatusMethodNotAllowed)
		return
	}

	education, err := h.educationRepo.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "GET error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if education == nil {
		http.Error(w, "Education not found", http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(education); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// Update handles HTTP requests to update an existing education resource.
//
// Behavior:
//   - Only supports the HTTP PUT method. Returns 405 Method Not Allowed for others.
//   - Expects a JSON request body representing domain.Education. Typical fields include:
//     Id, MainSchool, SchoolPeriods, Projects, Level, CreatedAt, UpdatedAt.
//   - Decodes the JSON payload and maps it to a domain.Education value.
//   - Validates the mapped education payload via Education.ValidatePayload(); returns 400 Bad Request on validation errors.
//   - Calls h.educationRepo.Update(ctx, education) to perform the persistent update.
//   - If the repository returns an error, responds with 500 Internal Server Error.
//   - If the repository returns nil (not found), responds with 404 Not Found.
//   - On success, encodes the updated education as JSON, sets Content-Type: application/json and responds with 200 OK.
//   - Ensures the request body is closed and propagates request context to the repository call.
//
// Notes:
// - This handler performs input validation before invoking the repository to avoid persisting invalid data.
// - All error responses include concise diagnostic messages and appropriate HTTP status codes.
//
// @Security ApiKeyAuth
// @Summary Update a education
// @Description Updates an existing education using the ID provided in the request body. Returns the updated education.
// @Tags education
// @Accept json
// @Produce json
// @Param education body domain.Education true "Education payload with ID"
// @Success 200 {object} domain.Education
// @Failure 400 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /education [put]
func (h *educationServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed: only PUT is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var updateReq domain.Education
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	if err := updateReq.ValidatePayload(); err != nil {
		http.Error(w, "Invalid education payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedEducation, err := h.educationRepo.Update(r.Context(), &updateReq)
	if err != nil {
		http.Error(w, "Failed to update education: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if updatedEducation == nil {
		http.Error(w, "Education not found", http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(updatedEducation); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
