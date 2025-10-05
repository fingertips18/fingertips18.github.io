package v1

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	mockRepo "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper to marshal JSON safely
func toJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

type projectHandlerTestFixture struct {
	t               *testing.T
	mockProjectRepo *mockRepo.MockProjectRepository
	projectHandler  ProjectHandler
}

func newProjectHandlerTestFixture(t *testing.T) *projectHandlerTestFixture {
	mockProjectRepo := new(mockRepo.MockProjectRepository)

	projectHandler := NewProjectServiceHandler(
		ProjectServiceConfig{
			projectRepo: mockProjectRepo,
		},
	)

	return &projectHandlerTestFixture{
		t:               t,
		mockProjectRepo: mockProjectRepo,
		projectHandler:  projectHandler,
	}
}

func TestProjectServiceHandler_Create(t *testing.T) {
	fixedID := "123-abc"

	validResp, _ := json.Marshal(domain.ProjectIDResponse{ID: fixedID})

	createReq := domain.CreateProject{
		Preview:     "preview.png",
		BlurHash:    "hash",
		Title:       "title",
		SubTitle:    "subtitle",
		Description: "desc",
		Stack:       []string{"go", "react"},
		Type:        domain.Web,
		Link:        "http://example.com",
	}
	validBody, _ := json.Marshal(createReq)

	type Given struct {
		method   string
		body     string
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
				method: http.MethodPost,
				body:   string(validBody),
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
		})
	}
}

func TestProjectServiceHandler_Get(t *testing.T) {
	fixedID := "123-abc"
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	validProject := &domain.Project{
		Id:          fixedID,
		Preview:     "preview.png",
		BlurHash:    "hash",
		Title:       "title",
		SubTitle:    "subtitle",
		Description: "desc",
		Stack:       []string{"go", "react"},
		Type:        domain.Web,
		Link:        "http://example.com",
		CreatedAt:   fixedTime,
		UpdatedAt:   fixedTime,
	}

	validBody, _ := json.Marshal(validProject)

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
		})
	}
}

func TestProjectServiceHandler_Update(t *testing.T) {
	fixedID := "123-abc"

	validProject := &domain.Project{
		Id:          fixedID,
		Preview:     "preview.png",
		BlurHash:    "hash",
		Title:       "title",
		SubTitle:    "subtitle",
		Description: "desc",
		Stack:       []string{"go", "react"},
		Type:        domain.Web,
		Link:        "http://example.com",
	}

	validBody, _ := json.Marshal(validProject)

	type Given struct {
		method   string
		body     string
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
				method: http.MethodPut,
				body:   string(validBody),
				mockRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						Update(mock.Anything, validProject).
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
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newProjectHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockProjectRepo)
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
		})
	}
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
		})
	}
}

func TestProjectServiceHandler_List(t *testing.T) {
	fixedTime := time.Now()
	validProjects := []domain.Project{
		{
			Id:          "p1",
			Preview:     "preview1",
			BlurHash:    "hash1",
			Title:       "title1",
			SubTitle:    "subtitle1",
			Description: "desc1",
			Stack:       []string{"go"},
			Type:        domain.Web,
			Link:        "http://example.com/1",
			CreatedAt:   fixedTime,
			UpdatedAt:   fixedTime,
		},
		{
			Id:          "p2",
			Preview:     "preview2",
			BlurHash:    "hash2",
			Title:       "title2",
			SubTitle:    "subtitle2",
			Description: "desc2",
			Stack:       []string{"react"},
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
				body: toJSON(validProjects),
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
				body: toJSON([]domain.Project{validProjects[0]}),
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
		})
	}
}
