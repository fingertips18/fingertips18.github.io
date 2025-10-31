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
	"github.com/jackc/pgx/v5"
)

type SkillHandler interface {
	http.Handler
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, id string)
	Update(w http.ResponseWriter, r *http.Request)
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

func (h *skillServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Normalize path by trimming trailing slash
	path := strings.TrimSuffix(r.URL.Path, "/")

	switch {
	// POST /skill
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

	case strings.HasPrefix(path, "/skill/"):
		id := strings.TrimPrefix(path, "/skill/")

		if id == "" {
			http.Error(w, "Skill ID is required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.Get(w, r, id)
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
// @Success 201 {object} IDResponse "Education ID"
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
