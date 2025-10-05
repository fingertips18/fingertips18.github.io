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

type ProjectHandler interface {
	http.Handler
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, id string)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request, id string)
	List(w http.ResponseWriter, r *http.Request)
}

type ProjectServiceConfig struct {
	DatabaseAPI database.DatabaseAPI

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
func NewProjectServiceHandler(cfg ProjectServiceConfig) ProjectHandler {
	projectRepo := cfg.projectRepo
	if projectRepo == nil {
		projectRepo = v1.NewProjectRepository(
			v1.ProjectRepositoryConfig{
				DatabaseAPI:  cfg.DatabaseAPI,
				ProjectTable: "Project",
			},
		)
	}

	return &projectServiceHandler{
		projectRepo: projectRepo,
	}
}

// ServeHTTP handles HTTP requests for project-related endpoints.
//
// It supports the following routes:
//   - GET    /projects         : List all projects.
//   - POST   /project          : Create a new project.
//   - PUT    /project          : Update an existing project.
//   - GET    /project/{id}     : Retrieve a project by its ID.
//   - DELETE /project/{id}     : Delete a project by its ID.
//
// For unsupported methods or unknown routes, it responds with appropriate HTTP error codes.
func (h *projectServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Normalize path by trimming trailing slash
	path := strings.TrimSuffix(r.URL.Path, "/")

	switch {
	// GET /projects
	case path == "/projects":
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.List(w, r)
		return

	// POST / PUT /project
	case path == "/project":
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodPut:
			h.Update(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return

	// GET / DELETE /project/{id}
	case strings.HasPrefix(path, "/project/"):
		id := strings.TrimPrefix(path, "/project/")

		if r.Method != http.MethodGet && r.Method != http.MethodDelete {
			return
		}

		if id == "" {
			http.Error(w, "Project ID is required", http.StatusBadRequest)
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
		return

	// Unknown route
	default:
		http.NotFound(w, r)
		return
	}
}

// Create handles HTTP POST requests to create a new project.
// It expects a JSON payload in the request body representing a project.
// On success, it responds with a JSON object containing the new project's ID and a status message.
// If the request method is not POST, the JSON is invalid, or project creation fails, it responds with an appropriate HTTP error.
//
// @Security ApiKeyAuth
// @Summary Create a project
// @Description Creates a new project from the provided JSON payload. Returns the created project with an assigned ID.
// @Tags project
// @Accept json
// @Produce json
// @Param project body domain.CreateProject true "Project payload"
// @Success 201 {string} domain.ProjectIDResponse "Project ID"
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /project [post]
func (h *projectServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed: only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var project domain.CreateProject
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

	resp := domain.ProjectIDResponse{ID: id}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

// Get handles HTTP GET requests for retrieving a project by its ID.
// It expects a JSON body containing the "id" field. If the request method is not GET,
// it responds with "Method not allowed". On success, it returns the project data in JSON format.
// If there is an error decoding the request, retrieving the project, or encoding the response,
// it responds with the appropriate HTTP error status and message.
//
// @Security ApiKeyAuth
// @Summary Get a project by ID
// @Description Retrieves the details of a specific project using its unique ID.
// @Tags project
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} domain.Project
// @Failure 400 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /project/{id} [get]
func (h *projectServiceHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed: only GET is supported", http.StatusMethodNotAllowed)
		return
	}

	project, err := h.projectRepo.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "GET error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if project == nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(project); err != nil {
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
//
// @Security ApiKeyAuth
// @Summary Update a project
// @Description Updates an existing project using the ID provided in the request body. Returns the updated project.
// @Tags project
// @Accept json
// @Produce json
// @Param project body domain.Project true "Project payload with ID"
// @Success 200 {object} domain.Project
// @Failure 400 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /project [put]
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
			Id:          project.Id,
			Preview:     project.Preview,
			BlurHash:    project.BlurHash,
			Title:       project.Title,
			SubTitle:    project.SubTitle,
			Description: project.Description,
			Stack:       project.Stack,
			Type:        domain.ProjectType(project.Type),
			Link:        project.Link,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		},
	)
	if err != nil {
		http.Error(w, "Failed to update project: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if updatedProject == nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(updatedProject); err != nil {
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
//
// @Security ApiKeyAuth
// @Summary Delete a project
// @Description Deletes an existing project by its unique ID provided in the path.
// @Tags project
// @Param id path string true "Project ID"
// @Success 204 "No Content"
// @Failure 400 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /project/{id} [delete]
func (h *projectServiceHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed: only DELETE is supported", http.StatusMethodNotAllowed)
		return
	}

	err := h.projectRepo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Project not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to delete project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Successful delete → 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// List handles HTTP GET requests to retrieve a list of projects based on the provided filter criteria.
// It expects a JSON-encoded ProjectFilter in the request body, decodes it, and queries the project repository.
// On success, it responds with a JSON object containing the list of projects and a status message.
// If the request method is not GET, the JSON is invalid, or an error occurs during processing, it returns an appropriate HTTP error response.
//
// @Security ApiKeyAuth
// @Summary List projects
// @Description Retrieves a paginated list of projects with optional filtering and sorting.
// @Tags project
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Number of items per page (default 10)"
// @Param sort_by query string false "Field to sort by" Enums(created_at, updated_at)
// @Param sort_ascending query bool false "Sort ascending order"
// @Param type query string false "Filter by project type" Enums(web, mobile, game)
// @Success 200 {array} domain.Project
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /projects [get]
func (h *projectServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed: only GET is supported", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	typeStr := q.Get("type")
	var projectType *domain.ProjectType
	if typeStr != "" {
		t := domain.ProjectType(typeStr)
		switch t {
		case domain.Web, domain.Mobile, domain.Game:
			projectType = &t
		default:
			http.Error(w, "invalid project type", http.StatusBadRequest)
			return
		}
	}

	sortBy, err := utils.GetQuerySortBy(q, "sort_by")
	if err != nil {
		http.Error(w, "invalid sort by", http.StatusBadRequest)
		return
	}

	filter := domain.ProjectFilter{
		Page:          utils.GetQueryInt32(q, "page", 1),
		PageSize:      utils.GetQueryInt32(q, "page_size", 20),
		SortBy:        sortBy,
		SortAscending: utils.GetQueryBool(q, "sort_ascending", false),
		Type:          projectType,
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

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(projects); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
