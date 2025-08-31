package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	v1 "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1"
)

type ProjectService interface {
	http.Handler
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
}

type ProjectServiceConfig struct {
	ConnectionString string

	projectRepo v1.ProjectRepository
}

type projectServiceHandler struct {
	projectRepo v1.ProjectRepository
}

// NewProjectServiceHandler creates and returns a new instance of ProjectService.
// It accepts a ProjectServiceConfig, which may include a custom project repository.
// If no repository is provided in the config, it initializes a default ProjectRepository
// using the provided connection string and a default table name.
// Returns a ProjectService implementation.
func NewProjectServiceHandler(cfg ProjectServiceConfig) ProjectService {
	projectRepo := cfg.projectRepo
	if projectRepo == nil {
		projectRepo = v1.NewProjectRepository(
			v1.ProjectRepositoryConfig{
				ConnectionString: cfg.ConnectionString,
				ProjectTable:     "Project",
			},
		)
	}

	return &projectServiceHandler{
		projectRepo: projectRepo,
	}
}

// ServeHTTP routes incoming HTTP requests for the project service based on the URL path.
// It supports the following endpoints:
//   - /create: Handles project creation.
//   - /get: Retrieves a project.
//   - /update: Updates an existing project.
//   - /delete: Deletes a project.
//   - /list: Lists all projects.
//
// For any other path, it responds with a 404 Not Found.
func (h *projectServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/project")

	switch path {
	case "/create":
		h.Create(w, r)
	case "/get":
		h.Get(w, r)
	case "/update":
		h.Update(w, r)
	case "/delete":
		h.Delete(w, r)
	case "/list":
		h.List(w, r)
	default:
		http.NotFound(w, r)
	}
}

// Create handles HTTP POST requests to create a new project.
// It expects a JSON payload in the request body representing a project.
// On success, it responds with a JSON object containing the new project's ID and a status message.
// If the request method is not POST, the JSON is invalid, or project creation fails, it responds with an appropriate HTTP error.
func (h *projectServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed: only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var project domain.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	id, err := h.projectRepo.Create(
		r.Context(),
		domain.Project{
			Preview:     project.Preview,
			BlurHash:    project.BlurHash,
			Title:       project.Title,
			SubTitle:    project.SubTitle,
			Description: project.Description,
			Stack:       project.Stack,
			Type:        domain.ProjectType(project.Type),
			Link:        project.Link,
		},
	)
	if err != nil {
		http.Error(w, "Failed to create project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"id": id, "status": "ok"}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// Get handles HTTP GET requests for retrieving a project by its ID.
// It expects a JSON body containing the "id" field. If the request method is not GET,
// it responds with "Method not allowed". On success, it returns the project data in JSON format.
// If there is an error decoding the request, retrieving the project, or encoding the response,
// it responds with the appropriate HTTP error status and message.
func (h *projectServiceHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed: only GET is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	project, err := h.projectRepo.Get(r.Context(), req.ID)
	if err != nil {
		http.Error(w, "Failed to get project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"project": project,
		"status":  "ok",
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// Update handles HTTP PUT requests to update an existing project.
// It expects a JSON-encoded project in the request body, decodes it,
// and attempts to update the project in the repository. If successful,
// it responds with the updated project and a status message in JSON format.
// Returns appropriate HTTP error responses for invalid methods, bad JSON,
// update failures, or response encoding errors.
func (h *projectServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed: only PUT is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var project domain.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	updatedProject, err := h.projectRepo.Update(
		r.Context(),
		domain.Project{
			Preview:     project.Preview,
			BlurHash:    project.BlurHash,
			Title:       project.Title,
			SubTitle:    project.SubTitle,
			Description: project.Description,
			Stack:       project.Stack,
			Type:        domain.ProjectType(project.Type),
			Link:        project.Link,
		},
	)
	if err != nil {
		http.Error(w, "Failed to update project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"updatedProject": updatedProject,
		"status":         "ok",
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// Delete handles HTTP DELETE requests to remove a project by its ID.
// It expects a JSON body containing the "id" field. If the request method is not DELETE,
// it responds with "405 Method Not Allowed". On successful deletion, it returns a JSON
// response with status "ok". If an error occurs during decoding, deletion, or response
// encoding, it responds with the appropriate HTTP error status and message.
func (h *projectServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed: only DELETE is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	err := h.projectRepo.Delete(r.Context(), req.ID)
	if err != nil {
		http.Error(w, "Failed to delete project: "+err.Error(), http.StatusInternalServerError)
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

// List handles HTTP GET requests to retrieve a list of projects based on the provided filter criteria.
// It expects a JSON-encoded ProjectFilter in the request body, decodes it, and queries the project repository.
// On success, it responds with a JSON object containing the list of projects and a status message.
// If the request method is not GET, the JSON is invalid, or an error occurs during processing, it returns an appropriate HTTP error response.
func (h *projectServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed: only GET is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var filter domain.ProjectFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	projects, err := h.projectRepo.List(
		r.Context(),
		domain.ProjectFilter{
			Page:          filter.Page,
			PageSize:      filter.PageSize,
			SortBy:        filter.SortBy,
			SortAscending: filter.SortAscending,
			Type:          filter.Type,
		},
	)
	if err != nil {
		http.Error(w, "Failed to list project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"projects": projects,
		"status":   "ok",
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
