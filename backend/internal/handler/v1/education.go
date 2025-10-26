package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	v1 "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/utils"
	"github.com/jackc/pgx/v5"
)

type EducationHandler interface {
	http.Handler
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, id string)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request, id string)
	List(w http.ResponseWriter, r *http.Request)
}

type EducationServiceConfig struct {
	DatabaseAPI database.DatabaseAPI
	ProjectRepo v1.ProjectRepository

	educationRepo v1.EducationRepository
}

type educationServiceHandler struct {
	educationRepo v1.EducationRepository
	projectRepo   v1.ProjectRepository
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

	projectRepo := cfg.ProjectRepo
	if projectRepo == nil {
		projectRepo = v1.NewProjectRepository(
			v1.ProjectRepositoryConfig{
				DatabaseAPI:  cfg.DatabaseAPI,
				ProjectTable: "Project",
			},
		)
	}

	return &educationServiceHandler{
		educationRepo: educationRepo,
		projectRepo:   projectRepo,
	}
}

func (h *educationServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Normalize path by trimming trailing slash
	path := strings.TrimSuffix(r.URL.Path, "/")

	switch {
	// GET /educations
	case path == "/educations":
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.List(w, r)
		return

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

// Create handles HTTP POST requests to create a new education record.
// It only supports the POST method; other methods receive a 405 Method Not Allowed.
// The handler decodes a JSON request body into a CreateEducationRequest and
// defers closing the request body. It maps the decoded payload to a domain.Education
// (populating MainSchool, SchoolPeriods, Projects and Level) and calls
// h.educationRepo.Create with the request context. On success it returns a JSON
// body containing the newly created ID (IDResponse) with
// Content-Type "application/json" and HTTP status 201 Created. If JSON decoding
// fails the handler responds with 400 Bad Request; if creation or response
// encoding fails it responds with 500 Internal Server Error.
//
// @Security ApiKeyAuth
// @Summary Create an education
// @Description Creates a new education from the provided JSON payload. Returns the created education with an assigned ID.
// @Tags education
// @Accept json
// @Produce json
// @Param education body CreateEducationRequest true "Education payload"
// @Success 201 {object} IDResponse "Education ID"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /education [post]
func (h *educationServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed: only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var createReq CreateEducationRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	// Map school periods
	schoolPeriods := make([]domain.SchoolPeriod, len(createReq.SchoolPeriods))
	for i, sp := range createReq.SchoolPeriods {
		schoolPeriods[i] = domain.SchoolPeriod{
			Link:        sp.Link,
			Name:        sp.Name,
			Description: sp.Description,
			Logo:        sp.Logo,
			BlurHash:    sp.BlurHash,
			Honor:       sp.Honor,
			StartDate:   sp.StartDate,
			EndDate:     sp.EndDate,
		}
	}

	// Map to Education
	education := &domain.Education{
		MainSchool: domain.SchoolPeriod{
			Link:        createReq.MainSchool.Link,
			Name:        createReq.MainSchool.Name,
			Description: createReq.MainSchool.Description,
			Logo:        createReq.MainSchool.Logo,
			BlurHash:    createReq.MainSchool.BlurHash,
			Honor:       createReq.MainSchool.Honor,
			StartDate:   createReq.MainSchool.StartDate,
			EndDate:     createReq.MainSchool.EndDate,
		},
		SchoolPeriods: schoolPeriods,
		Level:         domain.EducationLevel(createReq.Level),
	}

	// Validate before calling repository
	if err := education.ValidatePayload(); err != nil {
		http.Error(w, "Invalid education payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.educationRepo.Create(r.Context(), education)

	if err != nil {
		http.Error(w, "Failed to create education: "+err.Error(), http.StatusInternalServerError)
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

// Get handles HTTP GET requests for an education resource identified by id.
// It enforces the GET method and returns StatusMethodNotAllowed for other HTTP methods.
// The handler uses the request context to retrieve the education entity from the repository.
// If the repository returns an error, Get responds with StatusInternalServerError.
// If the requested education resource is not found, Get responds with StatusNotFound.
// On success it encodes the education entity as JSON, sets Content-Type to application/json,
// and writes the payload with StatusOK. Any encoding or write error results in an internal server error response.
//
// @Security ApiKeyAuth
// @Summary Get an education by ID
// @Description Retrieves the details of a specific education using its unique ID.
// @Tags education
// @Accept json
// @Produce json
// @Param id path string true "Education ID"
// @Success 200 {object} EducationDTO
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /education/{id} [get]
func (h *educationServiceHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed: only GET is supported", http.StatusMethodNotAllowed)
		return
	}

	educationRes, err := h.educationRepo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Education not found", http.StatusNotFound)
			return
		}
		http.Error(w, "GET error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if educationRes == nil {
		http.Error(w, "Education not found", http.StatusNotFound)
		return
	}

	// Fetch related projects
	projects, err := h.projectRepo.ListByEducationID(r.Context(), id)
	if err != nil {
		log.Printf("Failed to list projects by education ID %s: %v", id, err)
		// Projects are optional, so we can return empty array
		projects = []domain.Project{}
	}

	// Convert to DTOs
	projectDTOs := make([]ProjectDTO, len(projects))
	for i, p := range projects {
		projectDTOs[i] = ProjectDTO{
			Id:          p.Id,
			Preview:     p.Preview,
			BlurHash:    p.BlurHash,
			Title:       p.Title,
			SubTitle:    p.SubTitle,
			Description: p.Description,
			Stack:       p.Stack,
			Type:        string(p.Type),
			Link:        p.Link,
			EducationID: p.EducationID,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
	}

	education := EducationDTO{
		Id: educationRes.Id,
		MainSchool: SchoolPeriodDTO{
			Link:        educationRes.MainSchool.Link,
			Name:        educationRes.MainSchool.Name,
			Description: educationRes.MainSchool.Description,
			Logo:        educationRes.MainSchool.Logo,
			BlurHash:    educationRes.MainSchool.BlurHash,
			Honor:       educationRes.MainSchool.Honor,
			StartDate:   educationRes.MainSchool.StartDate,
			EndDate:     educationRes.MainSchool.EndDate,
		},
		SchoolPeriods: func() []SchoolPeriodDTO {
			periods := make([]SchoolPeriodDTO, len(educationRes.SchoolPeriods))
			for i, p := range educationRes.SchoolPeriods {
				periods[i] = SchoolPeriodDTO{
					Link:        p.Link,
					Name:        p.Name,
					Description: p.Description,
					Logo:        p.Logo,
					BlurHash:    p.BlurHash,
					Honor:       p.Honor,
					StartDate:   p.StartDate,
					EndDate:     p.EndDate,
				}
			}
			return periods
		}(),
		Projects:  projectDTOs,
		Level:     string(educationRes.Level),
		CreatedAt: educationRes.CreatedAt,
		UpdatedAt: educationRes.UpdatedAt,
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
//   - Expects a JSON request body representing EducationDTO. Typical fields include:
//     Id, MainSchool, SchoolPeriods, Projects, Level, CreatedAt, UpdatedAt.
//   - Decodes the JSON payload and maps it to a EducationDTO value.
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
// @Summary Update an education
// @Description Updates an existing education using the ID provided in the request body. Returns the updated education.
// @Tags education
// @Accept json
// @Produce json
// @Param education body UpdateEducationRequest true "Education payload with ID"
// @Success 200 {object} UpdateEducationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /education [put]
func (h *educationServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed: only PUT is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var updateReq UpdateEducationRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	// Map school periods
	schoolPeriods := make([]domain.SchoolPeriod, len(updateReq.SchoolPeriods))
	for i, sp := range updateReq.SchoolPeriods {
		schoolPeriods[i] = domain.SchoolPeriod{
			Link:        sp.Link,
			Name:        sp.Name,
			Description: sp.Description,
			Logo:        sp.Logo,
			BlurHash:    sp.BlurHash,
			Honor:       sp.Honor,
			StartDate:   sp.StartDate,
			EndDate:     sp.EndDate,
		}
	}

	education := domain.Education{
		Id: updateReq.Id,
		MainSchool: domain.SchoolPeriod{
			Link:        updateReq.MainSchool.Link,
			Name:        updateReq.MainSchool.Name,
			Description: updateReq.MainSchool.Description,
			Logo:        updateReq.MainSchool.Logo,
			BlurHash:    updateReq.MainSchool.BlurHash,
			Honor:       updateReq.MainSchool.Honor,
			StartDate:   updateReq.MainSchool.StartDate,
			EndDate:     updateReq.MainSchool.EndDate,
		},
		SchoolPeriods: schoolPeriods,
		Level:         domain.EducationLevel(updateReq.Level),
		UpdatedAt:     time.Now(),
	}

	if err := education.ValidatePayload(); err != nil {
		http.Error(w, "Invalid education payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedEducationRes, err := h.educationRepo.Update(r.Context(), &education)
	if err != nil {
		http.Error(w, "Failed to update education: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if updatedEducationRes == nil {
		http.Error(w, "Education not found", http.StatusNotFound)
		return
	}

	updatedEducation := UpdateEducationRequest{
		Id: updatedEducationRes.Id,
		MainSchool: SchoolPeriodDTO{
			Link:        updatedEducationRes.MainSchool.Link,
			Name:        updatedEducationRes.MainSchool.Name,
			Description: updatedEducationRes.MainSchool.Description,
			Logo:        updatedEducationRes.MainSchool.Logo,
			BlurHash:    updatedEducationRes.MainSchool.BlurHash,
			Honor:       updatedEducationRes.MainSchool.Honor,
			StartDate:   updatedEducationRes.MainSchool.StartDate,
			EndDate:     updatedEducationRes.MainSchool.EndDate,
		},
		SchoolPeriods: func() []SchoolPeriodDTO {
			periods := make([]SchoolPeriodDTO, len(updatedEducationRes.SchoolPeriods))
			for i, p := range updatedEducationRes.SchoolPeriods {
				periods[i] = SchoolPeriodDTO{
					Link:        p.Link,
					Name:        p.Name,
					Description: p.Description,
					Logo:        p.Logo,
					BlurHash:    p.BlurHash,
					Honor:       p.Honor,
					StartDate:   p.StartDate,
					EndDate:     p.EndDate,
				}
			}
			return periods
		}(),
		Level: string(updatedEducationRes.Level),
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

// Delete handles HTTP DELETE requests to remove an education resource identified by id.
// It enforces the DELETE method (returns 405 Method Not Allowed for other methods).
// The handler delegates deletion to the education repository using the request context.
// If the repository reports no matching row, Delete responds with 404 Not Found.
// For other repository errors it responds with 500 Internal Server Error and an error message.
// On successful deletion it writes a 204 No Content response with no body.
//
// @Security ApiKeyAuth
// @Summary Delete an education
// @Description Deletes an existing education by its unique ID provided in the path.
// @Tags education
// @Param id path string true "Education ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /education/{id} [delete]
func (h *educationServiceHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed: only DELETE is supported", http.StatusMethodNotAllowed)
		return
	}

	err := h.educationRepo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Education not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to delete education: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Successful delete → 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// List handles HTTP GET requests to list education records.
// It accepts query parameters:
//   - "page" (int, default 1)
//   - "page_size" (int, default 20)
//   - "sort_by" (validated by utils.GetQuerySortBy)
//   - "sort_ascending" (bool, default false)
//
// If the request method is not GET the handler responds with 405 Method Not Allowed.
// If "sort_by" is invalid the handler responds with 400 Bad Request.
// The handler constructs a EducationFilterRequest from the parsed parameters, calls
// h.educationRepo.List with the request context, and returns the result as JSON with
// Content-Type "application/json" and HTTP 200 on success. Repository or encoding errors
// result in a 500 Internal Server Error response.
//
// @Security ApiKeyAuth
// @Summary List educations
// @Description Retrieves a paginated list of educations with optional filtering and sorting.
// @Tags education
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Number of items per page (default 10)"
// @Param sort_by query string false "Field to sort by" Enums(created_at, updated_at)
// @Param sort_ascending query bool false "Sort ascending order"
// @Success 200 {array} EducationDTO
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /educations [get]
func (h *educationServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed: only GET is supported", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	sortBy, err := utils.GetQuerySortBy(q, "sort_by")
	if err != nil {
		http.Error(w, "invalid sort by", http.StatusBadRequest)
		return
	}

	filter := EducationFilterRequest{
		Page:          utils.GetQueryInt32(q, "page", 1),
		PageSize:      utils.GetQueryInt32(q, "page_size", 10),
		SortBy:        sortBy,
		SortAscending: utils.GetQueryBool(q, "sort_ascending", false),
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

	var sortByPtr *domain.SortBy
	if filter.SortBy != "" {
		sb := domain.SortBy(filter.SortBy)
		sortByPtr = &sb
	}
	domainFilter := domain.EducationFilter{
		Page:          filter.Page,
		PageSize:      filter.PageSize,
		SortBy:        sortByPtr,
		SortAscending: filter.SortAscending,
	}

	educationsRes, err := h.educationRepo.List(r.Context(), domainFilter)
	if err != nil {
		http.Error(w, "Failed to list educations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	educations := make([]EducationDTO, len(educationsRes))
	for i, e := range educationsRes {

		// Fetch related projects
		projects, err := h.projectRepo.ListByEducationID(r.Context(), e.Id)
		if err != nil {
			log.Printf("Failed to list projects by education ID %s: %v", e.Id, err)
			// Projects are optional, so we can return empty array
			projects = []domain.Project{}
		}

		// Convert to DTOs
		projectDTOs := make([]ProjectDTO, len(projects))
		for i, p := range projects {
			projectDTOs[i] = ProjectDTO{
				Id:          p.Id,
				Preview:     p.Preview,
				BlurHash:    p.BlurHash,
				Title:       p.Title,
				SubTitle:    p.SubTitle,
				Description: p.Description,
				Stack:       p.Stack,
				Type:        string(p.Type),
				Link:        p.Link,
				EducationID: p.EducationID,
				CreatedAt:   p.CreatedAt,
				UpdatedAt:   p.UpdatedAt,
			}
		}

		educations[i] = EducationDTO{
			Id: e.Id,
			MainSchool: SchoolPeriodDTO{
				Link:        e.MainSchool.Link,
				Name:        e.MainSchool.Name,
				Description: e.MainSchool.Description,
				Logo:        e.MainSchool.Logo,
				BlurHash:    e.MainSchool.BlurHash,
				Honor:       e.MainSchool.Honor,
				StartDate:   e.MainSchool.StartDate,
				EndDate:     e.MainSchool.EndDate,
			},
			SchoolPeriods: func() []SchoolPeriodDTO {
				periods := make([]SchoolPeriodDTO, len(e.SchoolPeriods))
				for j, p := range e.SchoolPeriods {
					periods[j] = SchoolPeriodDTO{
						Link:        p.Link,
						Name:        p.Name,
						Description: p.Description,
						Logo:        p.Logo,
						BlurHash:    p.BlurHash,
						Honor:       p.Honor,
						StartDate:   p.StartDate,
						EndDate:     p.EndDate,
					}
				}
				return periods
			}(),
			Projects:  projectDTOs,
			Level:     string(e.Level),
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(educations); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
