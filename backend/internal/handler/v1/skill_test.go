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

	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	mockRepo "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1/mocks"
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

			req := httptest.NewRequest(tt.given.method, "/skills", strings.NewReader(tt.given.body))
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
