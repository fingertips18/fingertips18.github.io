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
	mockRepo "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type skillHandlerTestFixture struct {
	t             *testing.T
	mockSkillRepo *mockRepo.MockSkillRepository
	skillHandler  SkillHandler
}

func newSkillHandlerTestFixture(t *testing.T) *skillHandlerTestFixture {
	mockSkillRepo := new(mockRepo.MockSkillRepository)
	skillHandler := NewSkillServiceHandler(
		SkillServiceConfig{
			skillRepo: mockSkillRepo,
		},
	)

	return &skillHandlerTestFixture{
		t:             t,
		mockSkillRepo: mockSkillRepo,
		skillHandler:  skillHandler,
	}
}

func TestSkillServiceHandler_Create(t *testing.T) {
	fixedID := "skill-123"

	validReq := CreateSkillRequest{
		Icon:     "code",
		HexColor: "#FF5733",
		Label:    "Go",
		Category: "backend",
	}
	validBody, _ := json.Marshal(validReq)
	validResp, _ := json.Marshal(IDResponse{Id: fixedID})

	type Given struct {
		method   string
		body     string
		mockRepo func(m *mockRepo.MockSkillRepository)
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
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Create(mock.Anything, mock.MatchedBy(func(s *domain.Skill) bool {
							return s.Label == "Go" && s.Icon == "code" &&
								s.HexColor == "#FF5733" && s.Category == domain.SkillCategory("backend")
						})).
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
				body:     `{"invalid":}`,
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid JSON in request body\n",
			},
		},
		"invalid skill payload": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					bad := validReq
					bad.Label = "" // invalid
					b, _ := json.Marshal(bad)
					return string(b)
				}(),
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid skill payload: label missing\n",
			},
		},
		"repo error": {
			given: Given{
				method: http.MethodPost,
				body:   string(validBody),
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Create(mock.Anything, mock.AnythingOfType("*domain.Skill")).
						Return("", errors.New("db failure"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to create skill: db failure\n",
			},
		},
		"unicode label": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					unicodeReq := validReq
					unicodeReq.Label = "プログラミング 🧠"
					b, _ := json.Marshal(unicodeReq)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Create(mock.Anything, mock.MatchedBy(func(s *domain.Skill) bool {
							return strings.Contains(s.Label, "プログラミング")
						})).
						Return(fixedID, nil)
				},
			},
			expected: Expected{
				code: http.StatusCreated,
				body: string(validResp),
			},
		},
		"large payload": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					large := validReq
					large.Label = strings.Repeat("A", 5000)
					b, _ := json.Marshal(large)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Create(mock.Anything, mock.AnythingOfType("*domain.Skill")).
						Return(fixedID, nil)
				},
			},
			expected: Expected{
				code: http.StatusCreated,
				body: string(validResp),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newSkillHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockSkillRepo)
			}

			req := httptest.NewRequest(tt.given.method, "/skill", strings.NewReader(tt.given.body))
			w := httptest.NewRecorder()

			f.skillHandler.Create(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)

			if strings.HasPrefix(tt.expected.body, "{") {
				assert.JSONEq(t, tt.expected.body, string(body))
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			f.mockSkillRepo.AssertExpectations(t)
		})
	}
}

func TestSkillServiceHandler_Create_Routing(t *testing.T) {
	fixedID := "skill-123"

	validReq := CreateSkillRequest{
		Icon:     "code",
		HexColor: "#FF5733",
		Label:    "Go",
		Category: "backend",
	}
	validBody, _ := json.Marshal(validReq)
	expectedResp, _ := json.Marshal(IDResponse{Id: fixedID})

	f := newSkillHandlerTestFixture(t)

	f.mockSkillRepo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(s *domain.Skill) bool {
			return s.Label == "Go" && s.Icon == "code"
		})).
		Return(fixedID, nil)

	req := httptest.NewRequest(http.MethodPost, "/skill", bytes.NewReader(validBody))
	w := httptest.NewRecorder()

	handlerHTTP, ok := f.skillHandler.(http.Handler)
	assert.True(t, ok, "skillHandler should implement http.Handler")
	handlerHTTP.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.JSONEq(t, string(expectedResp), string(body))

	f.mockSkillRepo.AssertExpectations(t)
}

func TestSkillServiceHandler_Get(t *testing.T) {
	fixedID := "skill-123"
	fixedTime := time.Date(2025, 10, 25, 21, 56, 4, 0, time.UTC)

	sampleSkill := &SkillDTO{
		Id:        fixedID,
		Icon:      "icon.png",
		HexColor:  "#FFFFFF",
		Label:     "Golang",
		Category:  "programming",
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}

	validResp, _ := json.Marshal(sampleSkill)

	type Given struct {
		method      string
		id          string
		mockSkillFn func(m *mockRepo.MockSkillRepository)
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
				mockSkillFn: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Get(mock.Anything, fixedID).
						Return(&domain.Skill{
							Id:        fixedID,
							Icon:      "icon.png",
							HexColor:  "#FFFFFF",
							Label:     "Golang",
							Category:  domain.SkillCategory("programming"),
							CreatedAt: fixedTime,
							UpdatedAt: fixedTime,
						}, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: string(validResp),
			},
		},
		"method not allowed": {
			given: Given{
				method:      http.MethodPost,
				id:          fixedID,
				mockSkillFn: nil,
			},
			expected: Expected{
				code: http.StatusMethodNotAllowed,
				body: "Method not allowed: only GET is supported\n",
			},
		},
		"skill not found (nil result)": {
			given: Given{
				method: http.MethodGet,
				id:     "nonexistent",
				mockSkillFn: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Get(mock.Anything, "nonexistent").
						Return(nil, nil)
				},
			},
			expected: Expected{
				code: http.StatusNotFound,
				body: "Skill not found\n",
			},
		},
		"skill not found (pgx.ErrNoRows)": {
			given: Given{
				method: http.MethodGet,
				id:     "missing-skill",
				mockSkillFn: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Get(mock.Anything, "missing-skill").
						Return(nil, pgx.ErrNoRows)
				},
			},
			expected: Expected{
				code: http.StatusNotFound,
				body: "Skill not found\n",
			},
		},
		"repository error": {
			given: Given{
				method: http.MethodGet,
				id:     fixedID,
				mockSkillFn: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Get(mock.Anything, fixedID).
						Return(nil, errors.New("database failure"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "GET error: database failure\n",
			},
		},
		"empty id": {
			given: Given{
				method: http.MethodGet,
				id:     "",
				mockSkillFn: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Get(mock.Anything, "").
						Return(nil, nil)
				},
			},
			expected: Expected{
				code: http.StatusNotFound,
				body: "Skill not found\n",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newSkillHandlerTestFixture(t)

			if tt.given.mockSkillFn != nil {
				tt.given.mockSkillFn(f.mockSkillRepo)
			}

			req := httptest.NewRequest(tt.given.method, "/skill/"+tt.given.id, nil)
			w := httptest.NewRecorder()

			f.skillHandler.Get(w, req, tt.given.id)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)

			assert.Equal(t, tt.expected.code, res.StatusCode)

			if strings.HasPrefix(tt.expected.body, "{") {
				assert.JSONEq(t, tt.expected.body, string(body))
			} else if strings.HasPrefix(tt.expected.body, "Failed to write response:") {
				assert.Contains(t, string(body), "Failed to write response:")
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			f.mockSkillRepo.AssertExpectations(t)
		})
	}
}

func TestSkillServiceHandler_Get_Routing(t *testing.T) {
	fixedID := "skill-123"

	sampleSkill := &domain.Skill{
		Id:        fixedID,
		Icon:      "icon.png",
		HexColor:  "#ABCDEF",
		Label:     "Python",
		Category:  domain.SkillCategory("language"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	dto := SkillDTO{
		Id:        sampleSkill.Id,
		Icon:      sampleSkill.Icon,
		HexColor:  sampleSkill.HexColor,
		Label:     sampleSkill.Label,
		Category:  string(sampleSkill.Category),
		CreatedAt: sampleSkill.CreatedAt,
		UpdatedAt: sampleSkill.UpdatedAt,
	}

	expectedResp, _ := json.Marshal(dto)

	f := newSkillHandlerTestFixture(t)

	f.mockSkillRepo.EXPECT().
		Get(mock.Anything, fixedID).
		Return(sampleSkill, nil)

	req := httptest.NewRequest(http.MethodGet, "/skill/"+fixedID, nil)
	w := httptest.NewRecorder()

	handlerInterface, ok := f.skillHandler.(http.Handler)
	assert.True(t, ok, "skillHandler should implement http.Handler")

	handlerInterface.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, string(expectedResp), string(body))

	f.mockSkillRepo.AssertExpectations(t)
}

func TestSkillServiceHandler_Update(t *testing.T) {
	fixedID := "skill-123"

	existingSkill := &domain.Skill{
		Id:        fixedID,
		Icon:      "icon-js",
		HexColor:  "#f7df1e",
		Label:     "JavaScript",
		Category:  domain.Backend,
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	// Request DTO (input to handler)
	requestDTO := UpdateSkillRequest{
		Id:       existingSkill.Id,
		Icon:     existingSkill.Icon,
		HexColor: existingSkill.HexColor,
		Label:    existingSkill.Label,
		Category: string(existingSkill.Category),
	}

	// Response DTO (output from handler)
	responseDTO := UpdateSkillResponse{
		Id:        existingSkill.Id,
		Icon:      existingSkill.Icon,
		HexColor:  existingSkill.HexColor,
		Label:     existingSkill.Label,
		Category:  string(existingSkill.Category),
		CreatedAt: existingSkill.CreatedAt,
		UpdatedAt: existingSkill.UpdatedAt,
	}

	validBody, _ := json.Marshal(requestDTO)
	validResp, _ := json.Marshal(responseDTO)

	type Given struct {
		method   string
		body     string
		mockRepo func(m *mockRepo.MockSkillRepository)
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
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Update(mock.Anything, mock.MatchedBy(func(s *domain.Skill) bool {
							return s.Id == fixedID && s.Label == "JavaScript"
						})).
						Return(existingSkill, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: string(validResp),
			},
		},
		"method not allowed": {
			given: Given{
				method:   http.MethodPost,
				body:     string(validBody),
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusMethodNotAllowed,
				body: "Method not allowed: only PUT is supported\n",
			},
		},
		"invalid JSON": {
			given: Given{
				method:   http.MethodPut,
				body:     `{"invalid":}`,
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid JSON in request body\n",
			},
		},
		"invalid payload (validation failure)": {
			given: Given{
				method: http.MethodPut,
				body: func() string {
					bad := requestDTO
					bad.Label = "" // assume ValidatePayload requires Label
					b, _ := json.Marshal(bad)
					return string(b)
				}(),
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid skill payload: label missing\n",
			},
		},
		"repository error": {
			given: Given{
				method: http.MethodPut,
				body:   string(validBody),
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Update(mock.Anything, mock.AnythingOfType("*domain.Skill")).
						Return(nil, errors.New("database failure"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to update skill: database failure\n",
			},
		},
		"skill not found": {
			given: Given{
				method: http.MethodPut,
				body:   string(validBody),
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Update(mock.Anything, mock.AnythingOfType("*domain.Skill")).
						Return(nil, nil)
				},
			},
			expected: Expected{
				code: http.StatusNotFound,
				body: "Skill not found\n",
			},
		},
		"unicode fields": {
			given: Given{
				method: http.MethodPut,
				body: func() string {
					unicode := requestDTO
					unicode.Label = "東京大学 🏫"
					b, _ := json.Marshal(unicode)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Update(mock.Anything, mock.MatchedBy(func(s *domain.Skill) bool {
							return strings.Contains(s.Label, "東京")
						})).
						Return(existingSkill, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: string(validResp),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newSkillHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockSkillRepo)
			}

			req := httptest.NewRequest(tt.given.method, "/skill", strings.NewReader(tt.given.body))
			w := httptest.NewRecorder()

			f.skillHandler.Update(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)

			if strings.HasPrefix(tt.expected.body, "{") {
				assert.JSONEq(t, tt.expected.body, string(body))
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			f.mockSkillRepo.AssertExpectations(t)
		})
	}
}

func TestSkillServiceHandler_Delete(t *testing.T) {
	fixedID := "skill-123"

	type Given struct {
		method   string
		id       string
		mockRepo func(m *mockRepo.MockSkillRepository)
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
				mockRepo: func(m *mockRepo.MockSkillRepository) {
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
		"method not allowed": {
			given: Given{
				method: http.MethodGet,
				id:     fixedID,
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					// No expectations; handler should block before calling repo
				},
			},
			expected: Expected{
				code: http.StatusMethodNotAllowed,
				body: "Method not allowed: only DELETE is supported\n",
			},
		},
		"skill not found (pgx.ErrNoRows)": {
			given: Given{
				method: http.MethodDelete,
				id:     "missing-id",
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Delete(mock.Anything, "missing-id").
						Return(pgx.ErrNoRows)
				},
			},
			expected: Expected{
				code: http.StatusNotFound,
				body: "Skill not found\n",
			},
		},
		"repository error": {
			given: Given{
				method: http.MethodDelete,
				id:     fixedID,
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Delete(mock.Anything, fixedID).
						Return(errors.New("database failure"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to delete skill: database failure\n",
			},
		},
		"empty id": {
			given: Given{
				method: http.MethodDelete,
				id:     "",
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						Delete(mock.Anything, "").
						Return(pgx.ErrNoRows)
				},
			},
			expected: Expected{
				code: http.StatusNotFound,
				body: "Skill not found\n",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newSkillHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockSkillRepo)
			}

			req := httptest.NewRequest(tt.given.method, "/skill/"+tt.given.id, nil)
			w := httptest.NewRecorder()

			f.skillHandler.Delete(w, req, tt.given.id)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)
			assert.Equal(t, tt.expected.body, string(body))

			f.mockSkillRepo.AssertExpectations(t)
		})
	}
}

func TestSkillServiceHandler_Delete_Routing(t *testing.T) {
	fixedID := "skill-123"

	f := newSkillHandlerTestFixture(t)

	// Setup mock expectation
	f.mockSkillRepo.EXPECT().
		Delete(mock.Anything, fixedID).
		Return(nil)

	// Create request and recorder
	req := httptest.NewRequest(http.MethodDelete, "/skill/"+fixedID, nil)
	w := httptest.NewRecorder()

	// Cast handler to http.Handler and call ServeHTTP
	handler, ok := f.skillHandler.(http.Handler)
	assert.True(t, ok, "skillHandler should implement http.Handler")
	handler.ServeHTTP(w, req)

	// Verify response
	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	assert.Equal(t, "", string(body))

	f.mockSkillRepo.AssertExpectations(t)
}

func TestSkillServiceHandler_List(t *testing.T) {
	fixedTime := time.Now()
	validSkills := []domain.Skill{
		{
			Id:        "s1",
			Icon:      "icon-js",
			HexColor:  "#f7df1e",
			Label:     "JavaScript",
			Category:  domain.Frontend,
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
		{
			Id:        "s2",
			Icon:      "icon-go",
			HexColor:  "#00ADD8",
			Label:     "Go",
			Category:  domain.Backend,
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
	}

	type Given struct {
		method   string
		query    string
		mockRepo func(m *mockRepo.MockSkillRepository)
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
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						List(mock.Anything, mock.AnythingOfType("domain.SkillFilter")).
						Return(validSkills, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: toJSON(validSkills),
			},
		},
		"success - with category filter": {
			given: Given{
				method: http.MethodGet,
				query:  "?category=frontend&page=1&page_size=5&sort_by=created_at&sort_ascending=true",
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						List(mock.Anything, mock.AnythingOfType("domain.SkillFilter")).
						Return([]domain.Skill{validSkills[0]}, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: toJSON([]domain.Skill{validSkills[0]}),
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
		"invalid skill category": {
			given: Given{
				method: http.MethodGet,
				query:  "?category=ai",
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "invalid skill category\n",
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
				mockRepo: func(m *mockRepo.MockSkillRepository) {
					m.EXPECT().
						List(mock.Anything, mock.AnythingOfType("domain.SkillFilter")).
						Return(nil, errors.New("db failure"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to list skills: db failure\n",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newSkillHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockSkillRepo)
			}

			req := httptest.NewRequest(tt.given.method, "/skills"+tt.given.query, nil)
			w := httptest.NewRecorder()

			f.skillHandler.(*skillServiceHandler).List(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)

			if strings.HasPrefix(tt.expected.body, "[") {
				assert.JSONEq(t, tt.expected.body, string(body))
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			f.mockSkillRepo.AssertExpectations(t)
		})
	}
}

func TestSkillServiceHandler_List_Routing(t *testing.T) {
	fixedTime := time.Now()
	validSkills := []domain.Skill{
		{
			Id:        "s1",
			Icon:      "icon-js",
			HexColor:  "#f7df1e",
			Label:     "JavaScript",
			Category:  domain.Frontend,
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
	}

	f := newSkillHandlerTestFixture(t)

	f.mockSkillRepo.EXPECT().
		List(mock.Anything, mock.AnythingOfType("domain.SkillFilter")).
		Return(validSkills, nil)

	req := httptest.NewRequest(http.MethodGet, "/skills", nil)
	w := httptest.NewRecorder()

	handler, ok := f.skillHandler.(http.Handler)
	assert.True(t, ok, "skillHandler should implement http.Handler")
	handler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, toJSON(validSkills), string(body))

	f.mockSkillRepo.AssertExpectations(t)
}
