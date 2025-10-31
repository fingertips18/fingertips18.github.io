package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	v1 "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/utils"
	"github.com/jackc/pgx/v5"
)

type SkillHandler interface {
	http.Handler
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, id string)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request, id string)
	List(w http.ResponseWriter, r *http.Request)
}

type SkillServiceConfig struct {
	DatabaseAPI database.DatabaseAPI

	skillRepo v1.SkillRepository
}

type skillServiceHandler struct {
	skillRepo v1.SkillRepository
}

// NewSkillServiceHandler returns a SkillHandler wired according to the provided
// SkillServiceConfig. If cfg.skillRepo is nil, a default v1.SkillRepository is
// created using cfg.DatabaseAPI and the "Skill" table name. The resulting handler
// uses the supplied or default repository to satisfy skill-related operations.
func NewSkillServiceHandler(cfg SkillServiceConfig) SkillHandler {
	skillRepo := cfg.skillRepo
	if skillRepo == nil {
		skillRepo = v1.NewSkillRepository(
			v1.SkillRepositoryConfig{
				DatabaseAPI: cfg.DatabaseAPI,
				SkillTable:  "Skill",
			},
		)
	}

	return &skillServiceHandler{
		skillRepo: skillRepo,
	}
}

// ServeHTTP dispatches HTTP requests for skillServiceHandler.
//
// It normalizes the request path by trimming a trailing slash, then routes
// requests to CRUD handlers based on exact path equality or prefix checks.
// The routing and behavior are as follows:
//
//   - GET  /skills         -> calls h.List(w, r)
//   - POST /skill          -> calls h.Create(w, r)
//   - PUT  /skill          -> calls h.Update(w, r)
//   - GET  /skill/{id}     -> calls h.Get(w, r, id)
//   - DELETE /skill/{id}   -> calls h.Delete(w, r, id)
//
// If a method is not allowed for a matched route, ServeHTTP responds with
// 405 Method Not Allowed. If a required skill ID segment is missing for a
// /skill/{id} route, it responds with 400 Bad Request. Unknown routes yield
// a 404 Not Found.
//
// Note: path matching is simple (case-sensitive string equality and prefix),
// and request body parsing/validation is delegated to the called handler
// methods, which are expected to write appropriate responses and status codes.
func (h *skillServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Normalize path by trimming trailing slash
	path := strings.TrimSuffix(r.URL.Path, "/")

	switch {
	// GET /skills
	case path == "/skills":
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.List(w, r)
		return

	// POST / PUT /skill
	case path == "/skill":
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodPut:
			h.Update(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return

	// GET / DELETE /skill/{id}
	case strings.HasPrefix(path, "/skill/"):
		id := strings.TrimPrefix(path, "/skill/")

		if id == "" {
			http.Error(w, "Skill ID is required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.Get(w, r, id)
		case http.MethodDelete:
			h.Delete(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	// Unknown route
	default:
		http.NotFound(w, r)
		return
	}
}

// Create handles HTTP POST requests to create a new skill.
//
// Expected behavior:
//
//   - Only the POST method is supported; other methods return 405 Method Not Allowed.
//
//   - The request body must be valid JSON matching CreateSkillRequest. Invalid JSON returns 400 Bad Request.
//
//   - The request JSON is decoded into a CreateSkillRequest and mapped to domain.Skill:
//     Icon, HexColor, Label are copied directly and Category is converted to domain.SkillCategory.
//
//   - The skill payload is validated via skill.ValidatePayload(); validation failures return 400 Bad Request.
//
//   - On success, the handler calls the repository's Create method with the request context.
//
//   - If repository creation succeeds, the handler responds with 201 Created and a JSON body containing an IDResponse:
//
//     {"Id": "<new-id>"}
//
// - Repository or encoding errors result in 500 Internal Server Error.
// - The handler sets Content-Type: application/json for successful responses and ensures the request body is closed.
//
// Context and error handling:
// - The repository call uses r.Context() so cancellations/deadlines propagate.
// - All error responses include an appropriate HTTP status code and a brief message describing the failure.
//
// @Security ApiKeyAuth
// @Summary Create a skill
// @Description Creates a new skill from the provided JSON payload. Returns the created skill with an assigned ID.
// @Tags skill
// @Accept json
// @Produce json
// @Param skill body CreateSkillRequest true "Skill payload"
// @Success 201 {object} IDResponse "Skill ID"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /skill [post]
func (h *skillServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed: only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var createReq CreateSkillRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	skill := domain.Skill{
		Icon:     createReq.Icon,
		HexColor: createReq.HexColor,
		Label:    createReq.Label,
		Category: domain.SkillCategory(createReq.Category),
	}

	// Validate before calling repository
	if err := skill.ValidatePayload(); err != nil {
		http.Error(w, "Invalid skill payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.skillRepo.Create(r.Context(), &skill)
	if err != nil {
		http.Error(w, "Failed to create skill: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := IDResponse{Id: id}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

// Get handles HTTP GET requests to retrieve a skill by its ID.
// It enforces the HTTP method to be GET and returns 405 Method Not Allowed for other methods.
// It calls the handler's skillRepo.Get with the request context and the provided id.
// If the repository reports no rows (or returns nil), Get responds with 404 Not Found.
// On repository errors it responds with 500 Internal Server Error and includes the error text.
// On success it maps the repository model to a SkillDTO, encodes it as JSON, sets Content-Type: application/json,
// and writes a 200 OK response with the encoded payload.
// If JSON encoding or response writing fails, Get responds with 500 Internal Server Error.
// The id parameter is expected to be the skill identifier (typically extracted from the request URL).
//
// @Security ApiKeyAuth
// @Summary Get a skill by ID
// @Description Retrieves the details of a specific skill using its unique ID.
// @Tags skill
// @Accept json
// @Produce json
// @Param id path string true "Skill ID"
// @Success 200 {object} SkillDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /skill/{id} [get]
func (h *skillServiceHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed: only GET is supported", http.StatusMethodNotAllowed)
		return
	}

	skillRes, err := h.skillRepo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Skill not found", http.StatusNotFound)
			return
		}
		http.Error(w, "GET error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if skillRes == nil {
		http.Error(w, "Skill not found", http.StatusNotFound)
		return
	}

	skill := SkillDTO{
		Id:        skillRes.Id,
		Icon:      skillRes.Icon,
		HexColor:  skillRes.HexColor,
		Label:     skillRes.Label,
		Category:  string(skillRes.Category),
		CreatedAt: skillRes.CreatedAt,
		UpdatedAt: skillRes.UpdatedAt,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(skill); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// Update handles HTTP PUT requests to update an existing skill.
//
// Behavior:
//   - Only accepts HTTP PUT; any other method results in 405 Method Not Allowed.
//   - Reads and closes the request body, expecting a JSON payload that matches UpdateSkillRequest
//     (including fields such as Id, Icon, HexColor, Label, Category).
//   - Decodes the JSON into a domain.Skill, converts the Category string to domain.SkillCategory,
//     and validates the resulting payload via skill.ValidatePayload(); validation failures return 400 Bad Request.
//   - Calls h.skillRepo.Update with the request context to persist the change.
//   - If the repository returns an error, responds with 500 Internal Server Error.
//   - If the repository returns nil (skill not found), responds with 404 Not Found.
//   - On success, encodes an UpdateSkillResponse containing the updated skill fields (Id, Icon, HexColor, Label,
//     Category, CreatedAt, UpdatedAt) as JSON, sets Content-Type: application/json, and returns 200 OK.
//
// Notes:
// - Uses json.Decoder/Encoder for request/response processing.
// - Errors include brief human-readable messages for client feedback.
//
// @Security ApiKeyAuth
// @Summary Update a skill
// @Description Updates an existing skill using the ID provided in the request body. Returns the updated skill.
// @Tags skill
// @Accept json
// @Produce json
// @Param skill body UpdateSkillRequest true "Skill payload with ID"
// @Success 200 {object} UpdateSkillResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /skill [put]
func (h *skillServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed: only PUT is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var updateReq UpdateSkillRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	skill := domain.Skill{
		Id:       updateReq.Id,
		Icon:     updateReq.Icon,
		HexColor: updateReq.HexColor,
		Label:    updateReq.Label,
		Category: domain.SkillCategory(updateReq.Category),
	}

	if err := skill.ValidatePayload(); err != nil {
		http.Error(w, "Invalid skill payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedSkillRes, err := h.skillRepo.Update(r.Context(), &skill)
	if err != nil {
		http.Error(w, "Failed to update skill: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if updatedSkillRes == nil {
		http.Error(w, "Skill not found", http.StatusNotFound)
		return
	}

	updatedSkill := UpdateSkillResponse{
		Id:        updatedSkillRes.Id,
		Icon:      updatedSkillRes.Icon,
		HexColor:  updatedSkillRes.HexColor,
		Label:     updatedSkillRes.Label,
		Category:  string(updatedSkillRes.Category),
		CreatedAt: updatedSkillRes.CreatedAt,
		UpdatedAt: updatedSkillRes.UpdatedAt,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(updatedSkill); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// Delete handles HTTP DELETE requests to remove a skill identified by id.
// It expects the request method to be DELETE; otherwise it responds with
// 405 Method Not Allowed. The request context is forwarded to the repository
// layer when attempting the delete.
//
// On success the handler responds with 204 No Content and no response body.
// If the repository returns pgx.ErrNoRows the handler responds with 404 Not Found.
// Any other repository error results in a 500 Internal Server Error and an
// error message written to the response.
//
// Parameters:
//   - w: http.ResponseWriter used to send the HTTP response.
//   - r: *http.Request representing the incoming HTTP request.
//   - id: string identifier of the skill to delete.
//
// @Security ApiKeyAuth
// @Summary Delete a skill
// @Description Deletes an existing skill by its unique ID provided in the path.
// @Tags skill
// @Param id path string true "Skill ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /skill/{id} [delete]
func (h *skillServiceHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed: only DELETE is supported", http.StatusMethodNotAllowed)
		return
	}

	err := h.skillRepo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Skill not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to delete skill: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Successful delete → 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// List handles HTTP GET requests to return a paginated list of skills.
//
// It only supports the GET method; other methods receive HTTP 405 Method Not Allowed.
//
// Query parameters:
//   - page (int, default: 1): 1-based page index. Values < 1 are clamped to 1.
//   - page_size (int, default: 10, max: 100): number of items per page. Values < 1 revert to the default; values > 100 are clamped to 100.
//   - sort_by (string): validated via utils.GetQuerySortBy; if invalid, the handler returns HTTP 400.
//   - sort_ascending (bool, default: false): whether to sort ascending.
//   - category (string): optional skill category. Supported categories are "Frontend", "Backend", "Tools", and "Others"; unknown categories result in HTTP 400.
//
// Behavior:
//   - Builds a domain.SkillFilter from the validated query parameters and calls the repository to obtain the skill list.
//   - On repository errors, responds with HTTP 500 and an error message.
//   - On successful retrieval, encodes the skills as JSON, sets Content-Type: application/json, and responds with HTTP 200.
//   - On JSON encoding errors, responds with HTTP 500.
//
// Notes:
//   - The handler returns concise HTTP error responses for invalid input (400), unsupported method (405), and internal failures (500).
//
// @Security ApiKeyAuth
// @Summary List skills
// @Description Retrieves a paginated list of skills with optional filtering and sorting.
// @Tags skill
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Number of items per page (default 10)"
// @Param sort_by query string false "Field to sort by" Enums(created_at, updated_at)
// @Param sort_ascending query bool false "Sort ascending order"
// @Param category query string false "Filter by skill category" Enums(frontend, backend, tools, others)
// @Success 200 {array} SkillDTO
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /skills [get]
func (h *skillServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed: only GET is supported", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	sortBy, err := utils.GetQuerySortBy(q, "sort_by")
	if err != nil {
		http.Error(w, "invalid sort_by: must be 'created_at' or 'updated_at'", http.StatusBadRequest)
		return
	}

	filter := SkillFilterRequest{
		Page:          utils.GetQueryInt32(q, "page", 1),
		PageSize:      utils.GetQueryInt32(q, "page_size", 10),
		SortBy:        sortBy,
		SortAscending: utils.GetQueryBool(q, "sort_ascending", false),
		Category:      q.Get("category"),
	}

	// Clamp page to minimum of 1
	if filter.Page < 1 {
		filter.Page = 1
	}

	// Clamp page_size to valid range
	const maxPageSize = 100
	if filter.PageSize < 1 {
		filter.PageSize = 10 // default
	} else if filter.PageSize > maxPageSize {
		filter.PageSize = maxPageSize
	}

	var skillCategory *domain.SkillCategory
	if filter.Category != "" {
		t := domain.SkillCategory(filter.Category)
		switch t {
		case domain.Frontend, domain.Backend, domain.Tools, domain.Others:
			skillCategory = &t
		default:
			http.Error(w, "invalid category: must be one of 'frontend', 'backend', 'tools', 'others'", http.StatusBadRequest)
			return
		}
	}

	var sortByPtr *domain.SortBy
	if filter.SortBy != "" {
		sb := domain.SortBy(filter.SortBy)
		sortByPtr = &sb
	}
	domainFilter := domain.SkillFilter{
		Page:          filter.Page,
		PageSize:      filter.PageSize,
		SortBy:        sortByPtr,
		SortAscending: filter.SortAscending,
		Category:      skillCategory,
	}

	skills, err := h.skillRepo.List(r.Context(), domainFilter)
	if err != nil {
		http.Error(w, "Failed to list skills: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(skills); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
