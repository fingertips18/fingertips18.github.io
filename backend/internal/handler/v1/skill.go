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

type SkillHandler interface {
	http.Handler
	Create(w http.ResponseWriter, r *http.Request)
}

type SkillServiceConfig struct {
	DatabaseAPI database.DatabaseAPI

	skillRepo v1.SkillRepository
}

type skillServiceHandler struct {
	skillRepo v1.SkillRepository
}

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
