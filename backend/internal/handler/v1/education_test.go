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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type educationHandlerTestFixture struct {
	t                 *testing.T
	mockEducationRepo *mockRepo.MockEducationRepository
	educationHandler  EducationHandler
}

func newEducationHandlerTestFixture(t *testing.T) *educationHandlerTestFixture {
	mockEducationRepo := new(mockRepo.MockEducationRepository)

	educationHandler := NewEducationServiceHandler(
		EducationServiceConfig{
			educationRepo: mockEducationRepo,
		},
	)

	return &educationHandlerTestFixture{
		t:                 t,
		mockEducationRepo: mockEducationRepo,
		educationHandler:  educationHandler,
	}
}

func TestEducationServiceHandler_Create(t *testing.T) {
	fixedID := "edu-123"

	validSchool := domain.SchoolPeriod{
		Name:        "Harvard University",
		Description: "Top-tier education",
		Logo:        "logo.png",
		BlurHash:    "hash123",
		StartDate:   time.Date(2015, 9, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
	}

	validCreateReq := domain.CreateEducation{
		MainSchool:    validSchool,
		SchoolPeriods: []domain.SchoolPeriod{validSchool},
		Projects:      nil,
		Level:         domain.College,
	}
	validBody, _ := json.Marshal(validCreateReq)
	validResp, _ := json.Marshal(domain.EducationIDResponse{ID: fixedID})

	type Given struct {
		method   string
		body     string
		mockRepo func(m *mockRepo.MockEducationRepository)
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
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Create(mock.Anything, mock.AnythingOfType("*domain.Education")).
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
		"empty body": {
			given: Given{
				method:   http.MethodPost,
				body:     "",
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
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Create(mock.Anything, mock.AnythingOfType("*domain.Education")).
						Return("", errors.New("db failure"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to create education: db failure\n",
			},
		},
		"invalid level": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					invalid := validCreateReq
					invalid.Level = "PhD" // invalid per isValid()
					b, _ := json.Marshal(invalid)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					// Still called since handler doesn’t validate input itself
					m.EXPECT().
						Create(mock.Anything, mock.AnythingOfType("*domain.Education")).
						Return(fixedID, nil)
				},
			},
			expected: Expected{
				code: http.StatusCreated,
				body: string(validResp),
			},
		},
		"missing required fields in school": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					bad := validCreateReq
					bad.MainSchool.Name = "" // invalid
					b, _ := json.Marshal(bad)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Create(mock.Anything, mock.AnythingOfType("*domain.Education")).
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
			f := newEducationHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockEducationRepo)
			}

			req := httptest.NewRequest(tt.given.method, "/education", strings.NewReader(tt.given.body))
			w := httptest.NewRecorder()

			f.educationHandler.Create(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)

			if strings.HasPrefix(tt.expected.body, "{") {
				assert.JSONEq(t, tt.expected.body, string(body))
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			f.mockEducationRepo.AssertExpectations(t)
		})
	}
}
