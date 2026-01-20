package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/handler/v1/dto"
	mockRepo "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1/mocks"
	metadata "github.com/fingertips18/fingertips18.github.io/backend/pkg/metadata/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	validBlurHash = "LEHV6nWB2yk8pyo0adR*.7kCMdnj"
)

// Helper to marshal JSON safely
func toJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

type projectHandlerTestFixture struct {
	t               *testing.T
	mockBlurHashAPI *metadata.MockBlurHashAPI
	mockProjectRepo *mockRepo.MockProjectRepository
	mockFileRepo    *mockRepo.MockFileRepository
	projectHandler  ProjectHandler
}

func newProjectHandlerTestFixture(t *testing.T) *projectHandlerTestFixture {
	mockBlurHashAPI := new(metadata.MockBlurHashAPI)
	mockProjectRepo := new(mockRepo.MockProjectRepository)
	mockFileRepo := new(mockRepo.MockFileRepository)
	projectHandler := NewProjectServiceHandler(
		ProjectServiceConfig{
			BlurHashAPI: mockBlurHashAPI,
			projectRepo: mockProjectRepo,
			fileRepo:    mockFileRepo,
		},
	)

	return &projectHandlerTestFixture{
		t:               t,
		mockBlurHashAPI: mockBlurHashAPI,
		mockProjectRepo: mockProjectRepo,
		mockFileRepo:    mockFileRepo,
		projectHandler:  projectHandler,
	}
}

func TestProjectServiceHandler_Create(t *testing.T) {
	fixedID := "123-abc"

	validResp, _ := json.Marshal(domain.ProjectIDResponse{ID: fixedID})

	createReq := dto.CreateProjectRequest{
		BlurHash:    validBlurHash,
		Title:       "title",
		Subtitle:    "subtitle",
		Description: "desc",
		Tags:        []string{"go", "react"},
		Type:        "web",
		Link:        "http://example.com",
	}
	validBody, _ := json.Marshal(createReq)

	invalidBlurHashReq := dto.CreateProjectRequest{
		BlurHash:    "invalid-hash",
		Title:       "title",
		Subtitle:    "subtitle",
		Description: "desc",
		Tags:        []string{"go", "react"},
		Type:        "web",
		Link:        "http://example.com",
	}
	invalidBlurHashBody, _ := json.Marshal(invalidBlurHashReq)

	type Given struct {
		method       string
		body         string
		mockBlurHash func(m *metadata.MockBlurHashAPI)
		mockRepo     func(m *mockRepo.MockProjectRepository)
		mockFileRepo func(m *mockRepo.MockFileRepository)
	}
	type Expected struct {
		code int
		body string
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"success": {
			given: Given{
				method: http.MethodPost,
				body:   string(validBody),
				mockBlurHash: func(m *metadata.MockBlurHashAPI) {
					m.EXPECT().IsValid(validBlurHash).Return(true).Once()
				},
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						Create(mock.Anything, mock.AnythingOfType("*domain.Project")).
						Return(fixedID, nil)
				},
			},
			expected: Expected{
				code: http.StatusCreated,
				body: string(validResp),
			},
		},
		"invalid blurhash": {
			given: Given{
				method: http.MethodPost,
				body:   string(invalidBlurHashBody),
				mockBlurHash: func(m *metadata.MockBlurHashAPI) {
					m.EXPECT().IsValid("invalid-hash").Return(false).Once()
				},
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid project payload: blurHash invalid\n",
			},
		},
		"invalid method": {
			given: Given{
				method:   http.MethodGet,
				body:     "",
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusMethodNotAllowed,
				body: "Method not allowed: only POST is supported\n",
			},
		},
		"invalid json": {
			given: Given{
				method:   http.MethodPost,
				body:     `{"invalid",}`,
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid JSON in request body\n",
			},
		},
		"repo error": {
			given: Given{
				method: http.MethodPost,
				body:   string(validBody),
				mockBlurHash: func(m *metadata.MockBlurHashAPI) {
					m.EXPECT().IsValid(validBlurHash).Return(true).Once()
				},
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						Create(mock.Anything, mock.AnythingOfType("*domain.Project")).
						Return("", errors.New("db failure"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to create project: db failure\n",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newProjectHandlerTestFixture(t)

			if tt.given.mockBlurHash != nil {
				tt.given.mockBlurHash(f.mockBlurHashAPI)
			}
			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockProjectRepo)
			}

			req := httptest.NewRequest(tt.given.method, "/project", strings.NewReader(tt.given.body))
			w := httptest.NewRecorder()

			f.projectHandler.Create(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)

			if strings.HasPrefix(tt.expected.body, "{") {
				assert.JSONEq(t, tt.expected.body, string(body))
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			f.mockProjectRepo.AssertExpectations(t)
			f.mockBlurHashAPI.AssertExpectations(t)
		})
	}
}

func TestProjectServiceHandler_Create_Routing(t *testing.T) {
	fixedID := "123-abc"

	f := newProjectHandlerTestFixture(t)

	// Setup valid input and expected output
	createReq := dto.CreateProjectRequest{
		BlurHash:    validBlurHash,
		Title:       "title",
		Subtitle:    "subtitle",
		Description: "desc",
		Tags:        []string{"go", "react"},
		Type:        "web",
		Link:        "http://example.com",
	}
	reqBody, _ := json.Marshal(createReq)
	expectedResp, _ := json.Marshal(domain.ProjectIDResponse{ID: fixedID})

	// Mock blurhash validation to return true for happy path
	f.mockBlurHashAPI.EXPECT().IsValid(validBlurHash).Return(true).Once()

	// Setup mock expectation
	f.mockProjectRepo.EXPECT().
		Create(mock.Anything, mock.AnythingOfType("*domain.Project")).
		Return(fixedID, nil)

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/project", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	// Ensure the handler implements http.Handler
	handler, ok := f.projectHandler.(http.Handler)
	assert.True(t, ok, "projectHandler should implement http.Handler")

	// Serve the request
	handler.ServeHTTP(w, req)

	// Validate response
	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.JSONEq(t, string(expectedResp), string(body))

	f.mockProjectRepo.AssertExpectations(t)
	f.mockBlurHashAPI.AssertExpectations(t)
}

func TestProjectServiceHandler_Get(t *testing.T) {
	fixedID := "123-abc"
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	validProject := &domain.Project{
		Id:          fixedID,
		BlurHash:    validBlurHash,
		Title:       "title",
		Subtitle:    "subtitle",
		Description: "desc",
		Tags:        []string{"go", "react"},
		Type:        domain.Web,
		Link:        "http://example.com",
		CreatedAt:   fixedTime,
		UpdatedAt:   fixedTime,
	}

	projectDTO := dto.ProjectDTO{
		ID:          fixedID,
		BlurHash:    validBlurHash,
		Title:       "title",
		Subtitle:    "subtitle",
		Description: "desc",
		Tags:        []string{"go", "react"},
		Type:        string(domain.Web),
		Link:        "http://example.com",
		Previews:    []dto.FileDTO{},
		CreatedAt:   fixedTime,
		UpdatedAt:   fixedTime,
	}

	validBody, _ := json.Marshal(projectDTO)

	type Given struct {
		method   string
		id       string
		mockRepo func(m *mockRepo.MockProjectRepository)
	}

	type Expected struct {
		code int
		body string
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"success": {
			given: Given{
				method: http.MethodGet,
				id:     fixedID,
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						Get(mock.Anything, fixedID).
						Return(validProject, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: string(validBody),
			},
		},
		"invalid method": {
			given: Given{
				method:   http.MethodPost,
				id:       fixedID,
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusMethodNotAllowed,
				body: "Method not allowed: only GET is supported\n",
			},
		},
		"repo error": {
			given: Given{
				method: http.MethodGet,
				id:     fixedID,
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						Get(mock.Anything, fixedID).
						Return(nil, errors.New("db failure"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "GET error: db failure\n",
			},
		},
		"not found": {
			given: Given{
				method: http.MethodGet,
				id:     "missing-id",
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						Get(mock.Anything, "missing-id").
						Return(nil, nil)
				},
			},
			expected: Expected{
				code: http.StatusNotFound,
				body: "Project not found\n",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newProjectHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockProjectRepo)
			}

			if name == "success" {
				f.mockFileRepo.EXPECT().
					FindByParent(mock.Anything, "project", tt.given.id, domain.Image).
					Return([]domain.File{}, nil)
			}

			req := httptest.NewRequest(tt.given.method, "/project/"+tt.given.id, nil)
			w := httptest.NewRecorder()

			f.projectHandler.(*projectServiceHandler).Get(w, req, tt.given.id)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)

			if strings.HasPrefix(tt.expected.body, "{") {
				assert.JSONEq(t, tt.expected.body, string(body))
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			f.mockProjectRepo.AssertExpectations(t)
			if name == "success" {
				f.mockFileRepo.AssertExpectations(t)
			}
		})
	}
}

func TestProjectServiceHandler_Get_Routing(t *testing.T) {
	fixedID := "123-abc"
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	validProject := &domain.Project{
		Id:          fixedID,
		BlurHash:    validBlurHash,
		Title:       "title",
		Subtitle:    "subtitle",
		Description: "desc",
		Tags:        []string{"go", "react"},
		Type:        domain.Web,
		Link:        "http://example.com",
		CreatedAt:   fixedTime,
		UpdatedAt:   fixedTime,
	}
	expectedDTO := dto.ProjectDTO{
		ID:          fixedID,
		BlurHash:    validBlurHash,
		Title:       "title",
		Subtitle:    "subtitle",
		Description: "desc",
		Tags:        []string{"go", "react"},
		Type:        string(domain.Web),
		Link:        "http://example.com",
		Previews:    []dto.FileDTO{},
		CreatedAt:   fixedTime,
		UpdatedAt:   fixedTime,
	}
	expectedBody, _ := json.Marshal(expectedDTO)

	f := newProjectHandlerTestFixture(t)

	// Setup mock expectation
	f.mockProjectRepo.EXPECT().
		Get(mock.Anything, fixedID).
		Return(validProject, nil)

	f.mockFileRepo.EXPECT().
		FindByParent(mock.Anything, "project", fixedID, domain.Image).
		Return([]domain.File{}, nil)

	// Create HTTP request and response recorder
	req := httptest.NewRequest(http.MethodGet, "/project/"+fixedID, nil)
	w := httptest.NewRecorder()

	// Ensure handler implements http.Handler
	handler, ok := f.projectHandler.(http.Handler)
	assert.True(t, ok, "projectHandler should implement http.Handler")

	// Serve the HTTP request
	handler.ServeHTTP(w, req)

	// Verify the response
	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, string(expectedBody), string(body))

	f.mockProjectRepo.AssertExpectations(t)
}

func TestProjectServiceHandler_Update(t *testing.T) {
	fixedID := "123-abc"

	validProject := &domain.Project{
		Id:          fixedID,
		BlurHash:    validBlurHash,
		Title:       "title",
		Subtitle:    "subtitle",
		Description: "desc",
		Tags:        []string{"go", "react"},
		Type:        domain.Web,
		Link:        "http://example.com",
	}

	validBody, _ := json.Marshal(validProject)

	invalidBlurHashReq := dto.CreateProjectRequest{
		BlurHash:    "invalid-hash",
		Title:       "title",
		Subtitle:    "subtitle",
		Description: "desc",
		Tags:        []string{"go", "react"},
		Type:        "web",
		Link:        "http://example.com",
	}
	invalidBlurHashBody, _ := json.Marshal(invalidBlurHashReq)

	type Given struct {
		method       string
		body         string
		mockBlurHash func(m *metadata.MockBlurHashAPI)
		mockRepo     func(m *mockRepo.MockProjectRepository)
		mockFileRepo func(m *mockRepo.MockFileRepository)
	}

	type Expected struct {
		code int
		body string
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"success": {
			given: Given{
				method: http.MethodPut,
				body:   string(validBody),
				mockBlurHash: func(m *metadata.MockBlurHashAPI) {
					m.EXPECT().IsValid(validBlurHash).Return(true).Once()
				},
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						Update(mock.Anything, validProject).
						Return(validProject, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: toJSON(dto.ProjectDTO{
					ID:          fixedID,
					BlurHash:    validBlurHash,
					Title:       "title",
					Subtitle:    "subtitle",
					Description: "desc",
					Tags:        []string{"go", "react"},
					Type:        string(domain.Web),
					Link:        "http://example.com",
					Previews:    []dto.FileDTO{},
					CreatedAt:   validProject.CreatedAt,
					UpdatedAt:   validProject.UpdatedAt,
				}),
			},
		},
		"invalid method": {
			given: Given{
				method: http.MethodPost,
				body:   string(validBody),
			},
			expected: Expected{
				code: http.StatusMethodNotAllowed,
				body: "Method not allowed: only PUT is supported\n",
			},
		},
		"invalid json": {
			given: Given{
				method: http.MethodPut,
				body:   `{"invalid",}`,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid JSON in request body\n",
			},
		},
		"repo error": {
			given: Given{
				method: http.MethodPut,
				body:   string(validBody),
				mockBlurHash: func(m *metadata.MockBlurHashAPI) {
					m.EXPECT().IsValid(validBlurHash).Return(true).Once()
				},
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						Update(mock.Anything, validProject).
						Return(nil, errors.New("db failure"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to update project: db failure\n",
			},
		},
		"not found": {
			given: Given{
				method: http.MethodPut,
				body:   string(validBody),
				mockBlurHash: func(m *metadata.MockBlurHashAPI) {
					m.EXPECT().IsValid(validBlurHash).Return(true).Once()
				},
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						Update(mock.Anything, validProject).
						Return(nil, nil)
				},
			},
			expected: Expected{
				code: http.StatusNotFound,
				body: "Project not found\n",
			},
		},
		"invalid blurhash": {
			given: Given{
				method: http.MethodPut,
				body:   string(invalidBlurHashBody),
				mockBlurHash: func(m *metadata.MockBlurHashAPI) {
					m.EXPECT().IsValid("invalid-hash").Return(false).Once()
				},
				mockRepo: nil, // Should fail before repo call
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid project payload: blurHash invalid\n",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newProjectHandlerTestFixture(t)

			if tt.given.mockBlurHash != nil {
				tt.given.mockBlurHash(f.mockBlurHashAPI)
			}

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockProjectRepo)
			}

			if name == "success" {
				f.mockFileRepo.EXPECT().
					FindByParent(mock.Anything, "project", fixedID, domain.Image).
					Return([]domain.File{}, nil)
			}

			req := httptest.NewRequest(tt.given.method, "/project", strings.NewReader(tt.given.body))
			w := httptest.NewRecorder()

			f.projectHandler.(*projectServiceHandler).Update(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)

			if strings.HasPrefix(tt.expected.body, "{") {
				assert.JSONEq(t, tt.expected.body, string(body))
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			f.mockProjectRepo.AssertExpectations(t)
			f.mockBlurHashAPI.AssertExpectations(t)

			if name == "success" {
				f.mockFileRepo.AssertExpectations(t)
			}
		})
	}
}

func TestProjectServiceHandler_Update_Routing(t *testing.T) {
	fixedID := "123-abc"
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	validProject := &domain.Project{
		Id:          fixedID,
		BlurHash:    validBlurHash,
		Title:       "title",
		Subtitle:    "subtitle",
		Description: "desc",
		Tags:        []string{"go", "react"},
		Type:        domain.Web,
		Link:        "http://example.com",
		CreatedAt:   fixedTime,
		UpdatedAt:   fixedTime,
	}

	reqBody, _ := json.Marshal(validProject)

	expectedResp, _ := json.Marshal(dto.ProjectDTO{
		ID:          fixedID,
		BlurHash:    validBlurHash,
		Title:       "title",
		Subtitle:    "subtitle",
		Description: "desc",
		Tags:        []string{"go", "react"},
		Type:        string(domain.Web),
		Link:        "http://example.com",
		Previews:    []dto.FileDTO{},
		CreatedAt:   fixedTime,
		UpdatedAt:   fixedTime,
	})

	f := newProjectHandlerTestFixture(t)

	// Mock blurhash validation to return true for happy path
	f.mockBlurHashAPI.EXPECT().IsValid(validBlurHash).Return(true).Once()

	// Mock projectRepo.Update call
	f.mockProjectRepo.EXPECT().
		Update(mock.Anything, validProject).
		Return(validProject, nil)

	// Mock fileRepo.FindByParent call to reload previews
	f.mockFileRepo.EXPECT().
		FindByParent(mock.Anything, "project", fixedID, domain.Image).
		Return([]domain.File{}, nil)

	// Create PUT request
	req := httptest.NewRequest(http.MethodPut, "/project", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	// Ensure handler implements http.Handler
	handler, ok := f.projectHandler.(http.Handler)
	assert.True(t, ok, "projectHandler should implement http.Handler")

	// Serve request
	handler.ServeHTTP(w, req)

	// Verify response
	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, string(expectedResp), string(body))

	f.mockProjectRepo.AssertExpectations(t)
	f.mockFileRepo.AssertExpectations(t)
}

func TestProjectServiceHandler_Delete(t *testing.T) {
	fixedID := "123-abc"

	type Given struct {
		method   string
		id       string
		mockRepo func(m *mockRepo.MockProjectRepository)
	}

	type Expected struct {
		code int
		body string
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"success": {
			given: Given{
				method: http.MethodDelete,
				id:     fixedID,
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						Delete(mock.Anything, fixedID).
						Return(nil)
				},
			},
			expected: Expected{
				code: http.StatusNoContent,
				body: "",
			},
		},
		"invalid method": {
			given: Given{
				method: http.MethodGet,
				id:     fixedID,
			},
			expected: Expected{
				code: http.StatusMethodNotAllowed,
				body: "Method not allowed: only DELETE is supported\n",
			},
		},
		"not found": {
			given: Given{
				method: http.MethodDelete,
				id:     fixedID,
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						Delete(mock.Anything, fixedID).
						Return(pgx.ErrNoRows)
				},
			},
			expected: Expected{
				code: http.StatusNotFound,
				body: "Project not found\n",
			},
		},
		"repo error": {
			given: Given{
				method: http.MethodDelete,
				id:     fixedID,
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						Delete(mock.Anything, fixedID).
						Return(errors.New("db failure"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to delete project: db failure\n",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newProjectHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockProjectRepo)
			}

			if name == "success" {
				f.mockFileRepo.EXPECT().
					DeleteByParent(mock.Anything, "project", tt.given.id).
					Return(nil)
			}

			req := httptest.NewRequest(tt.given.method, "/project/"+tt.given.id, nil)
			w := httptest.NewRecorder()

			f.projectHandler.(*projectServiceHandler).Delete(w, req, tt.given.id)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)

			if tt.expected.body != "" {
				assert.Equal(t, tt.expected.body, string(body))
			} else {
				assert.Empty(t, string(body))
			}

			f.mockProjectRepo.AssertExpectations(t)

			if name == "success" {
				f.mockFileRepo.AssertExpectations(t)
			}
		})
	}
}

func TestProjectServiceHandler_Delete_Routing(t *testing.T) {
	fixedID := "123-abc"

	f := newProjectHandlerTestFixture(t)

	f.mockFileRepo.EXPECT().
		DeleteByParent(mock.Anything, "project", fixedID).
		Return(nil)

	// Setup mock expectation
	f.mockProjectRepo.EXPECT().
		Delete(mock.Anything, fixedID).
		Return(nil)

	// Create DELETE request
	req := httptest.NewRequest(http.MethodDelete, "/project/"+fixedID, nil)
	w := httptest.NewRecorder()

	// Ensure the handler implements http.Handler
	handler, ok := f.projectHandler.(http.Handler)
	assert.True(t, ok, "projectHandler should implement http.Handler")

	// Serve the request
	handler.ServeHTTP(w, req)

	// Verify the response
	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	assert.Empty(t, string(body))

	// Verify the mock was called as expected
	f.mockFileRepo.AssertExpectations(t)
	f.mockProjectRepo.AssertExpectations(t)
}

func TestProjectServiceHandler_List(t *testing.T) {
	fixedTime := time.Now()
	validProjects := []domain.Project{
		{
			Id:          "p1",
			BlurHash:    "hash1",
			Title:       "title1",
			Subtitle:    "subtitle1",
			Description: "desc1",
			Tags:        []string{"go"},
			Type:        domain.Web,
			Link:        "http://example.com/1",
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
		},
		{
			Id:          "p2",
			BlurHash:    "hash2",
			Title:       "title2",
			Subtitle:    "subtitle2",
			Description: "desc2",
			Tags:        []string{"react"},
			Type:        domain.Mobile,
			Link:        "http://example.com/2",
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
		},
	}

	type Given struct {
		method   string
		query    string
		mockRepo func(m *mockRepo.MockProjectRepository)
	}
	type Expected struct {
		code int
		body string
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"success - no filters": {
			given: Given{
				method: http.MethodGet,
				query:  "",
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						List(mock.Anything, mock.AnythingOfType("domain.ProjectFilter")).
						Return(validProjects, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: toJSON([]dto.ProjectDTO{
					{
						ID:          "p1",
						BlurHash:    "hash1",
						Title:       "title1",
						Subtitle:    "subtitle1",
						Description: "desc1",
						Tags:        []string{"go"},
						Type:        string(domain.Web),
						Link:        "http://example.com/1",
						Previews:    []dto.FileDTO{},
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
					{
						ID:          "p2",
						BlurHash:    "hash2",
						Title:       "title2",
						Subtitle:    "subtitle2",
						Description: "desc2",
						Tags:        []string{"react"},
						Type:        string(domain.Mobile),
						Link:        "http://example.com/2",
						Previews:    []dto.FileDTO{},
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				}),
			},
		},
		"success - with type filter": {
			given: Given{
				method: http.MethodGet,
				query:  "?type=web&page=1&page_size=10&sort_by=created_at&sort_ascending=true",
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						List(mock.Anything, mock.AnythingOfType("domain.ProjectFilter")).
						Return([]domain.Project{validProjects[0]}, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: toJSON([]dto.ProjectDTO{
					{
						ID:          "p1",
						BlurHash:    "hash1",
						Title:       "title1",
						Subtitle:    "subtitle1",
						Description: "desc1",
						Tags:        []string{"go"},
						Type:        string(domain.Web),
						Link:        "http://example.com/1",
						Previews:    []dto.FileDTO{},
						CreatedAt:   fixedTime,
						UpdatedAt:   fixedTime,
					},
				}),
			},
		},
		"invalid method": {
			given: Given{
				method: http.MethodPost,
				query:  "",
			},
			expected: Expected{
				code: http.StatusMethodNotAllowed,
				body: "Method not allowed: only GET is supported\n",
			},
		},
		"invalid project type": {
			given: Given{
				method: http.MethodGet,
				query:  "?type=desktop",
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "invalid project type\n",
			},
		},
		"invalid sort by": {
			given: Given{
				method: http.MethodGet,
				query:  "?sort_by=invalid_field",
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "invalid sort by\n",
			},
		},
		"repo error": {
			given: Given{
				method: http.MethodGet,
				query:  "",
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						List(mock.Anything, mock.AnythingOfType("domain.ProjectFilter")).
						Return(nil, errors.New("db failure"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to list project: db failure\n",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newProjectHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockProjectRepo)
			}

			switch name {
			case "success - no filters":
				f.mockFileRepo.EXPECT().
					FindByParent(mock.Anything, "project", "p1", domain.Image).
					Return([]domain.File{}, nil)
				f.mockFileRepo.EXPECT().
					FindByParent(mock.Anything, "project", "p2", domain.Image).
					Return([]domain.File{}, nil)
			case "success - with type filter":
				f.mockFileRepo.EXPECT().
					FindByParent(mock.Anything, "project", "p1", domain.Image).
					Return([]domain.File{}, nil)
			}

			req := httptest.NewRequest(tt.given.method, "/projects"+tt.given.query, nil)
			w := httptest.NewRecorder()

			f.projectHandler.(*projectServiceHandler).List(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)

			if strings.HasPrefix(tt.expected.body, "[") {
				assert.JSONEq(t, tt.expected.body, string(body))
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			f.mockProjectRepo.AssertExpectations(t)
			if name == "success - no filters" || name == "success - with type filter" {
				f.mockFileRepo.AssertExpectations(t)
			}
		})
	}
}

func TestProjectServiceHandler_List_Routing(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	validProjects := []domain.Project{
		{
			Id:          "p1",
			BlurHash:    "hash1",
			Title:       "title1",
			Subtitle:    "subtitle1",
			Description: "desc1",
			Tags:        []string{"go"},
			Type:        domain.Web,
			Link:        "http://example.com/1",
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
		},
		{
			Id:          "p2",
			BlurHash:    "hash2",
			Title:       "title2",
			Subtitle:    "subtitle2",
			Description: "desc2",
			Tags:        []string{"react"},
			Type:        domain.Mobile,
			Link:        "http://example.com/2",
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
		},
	}

	// Expected response is dto.ProjectDTO array
	expectedBody, _ := json.Marshal([]dto.ProjectDTO{
		{
			ID:          "p1",
			BlurHash:    "hash1",
			Title:       "title1",
			Subtitle:    "subtitle1",
			Description: "desc1",
			Tags:        []string{"go"},
			Type:        string(domain.Web),
			Link:        "http://example.com/1",
			Previews:    []dto.FileDTO{},
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
		},
		{
			ID:          "p2",
			BlurHash:    "hash2",
			Title:       "title2",
			Subtitle:    "subtitle2",
			Description: "desc2",
			Tags:        []string{"react"},
			Type:        string(domain.Mobile),
			Link:        "http://example.com/2",
			Previews:    []dto.FileDTO{},
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
		},
	})

	f := newProjectHandlerTestFixture(t)

	// Mock the repository call
	f.mockProjectRepo.EXPECT().
		List(mock.Anything, mock.AnythingOfType("domain.ProjectFilter")).
		Return(validProjects, nil)

	// Mock fileRepo.FindByParent for each project
	f.mockFileRepo.EXPECT().
		FindByParent(mock.Anything, "project", "p1", domain.Image).
		Return([]domain.File{}, nil)
	f.mockFileRepo.EXPECT().
		FindByParent(mock.Anything, "project", "p2", domain.Image).
		Return([]domain.File{}, nil)

	// Create GET request to /projects
	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	w := httptest.NewRecorder()

	// Ensure the handler implements http.Handler
	handler, ok := f.projectHandler.(http.Handler)
	assert.True(t, ok, "projectHandler should implement http.Handler")

	// Serve through full HTTP routing
	handler.ServeHTTP(w, req)

	// Check response
	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, string(expectedBody), string(body))

	f.mockFileRepo.AssertExpectations(t)
	f.mockProjectRepo.AssertExpectations(t)
}
