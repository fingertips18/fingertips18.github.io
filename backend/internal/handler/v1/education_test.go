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
						Create(mock.Anything, mock.MatchedBy(func(edu *domain.Education) bool {
							return edu.Level == domain.College &&
								edu.MainSchool.Name == "Harvard University" &&
								len(edu.SchoolPeriods) == 1
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
				mockRepo: nil, // No mock needed - validation happens before repo call
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid education payload: level invalid = PhD\n",
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
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid education payload: main school name missing\n",
			},
		}, "very large payload": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					bigDesc := strings.Repeat("A", 10_000) // simulate large input
					large := validCreateReq
					large.MainSchool.Description = bigDesc
					b, _ := json.Marshal(large)
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
		"unicode characters in fields": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					unicodeReq := validCreateReq
					unicodeReq.MainSchool.Name = "Université de Montréal 🏫"
					unicodeReq.MainSchool.Description = "研究と教育 excellence"
					b, _ := json.Marshal(unicodeReq)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Create(mock.Anything, mock.MatchedBy(func(edu *domain.Education) bool {
							return strings.Contains(edu.MainSchool.Name, "Montréal")
						})).
						Return(fixedID, nil)
				},
			},
			expected: Expected{
				code: http.StatusCreated,
				body: string(validResp),
			},
		},
		"start date after end date": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					invalidDates := validCreateReq
					invalidDates.MainSchool.StartDate = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
					invalidDates.MainSchool.EndDate = time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
					b, _ := json.Marshal(invalidDates)
					return string(b)
				}(),
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid education payload: main school end date must be after start date\n",
			},
		},
		"multiple school periods": {
			given: Given{
				method: http.MethodPost,
				body: func() string {
					multi := validCreateReq
					secondSchool := domain.SchoolPeriod{
						Name:        "Massachusetts Institute of Technology",
						Description: "Exchange program in Computer Science",
						Logo:        "mit_logo.png",
						BlurHash:    "mitHash456",
						StartDate:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
						EndDate:     time.Date(2018, 6, 1, 0, 0, 0, 0, time.UTC),
					}
					multi.SchoolPeriods = []domain.SchoolPeriod{
						validSchool,
						secondSchool,
					}
					b, _ := json.Marshal(multi)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Create(mock.Anything, mock.MatchedBy(func(edu *domain.Education) bool {
							return len(edu.SchoolPeriods) == 2 &&
								edu.SchoolPeriods[1].Name == "Massachusetts Institute of Technology"
						})).
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

func TestEducationServiceHandler_Get(t *testing.T) {
	fixedID := "edu-123"

	sampleEducation := &domain.Education{
		Id: fixedID,
		MainSchool: domain.SchoolPeriod{
			Name:        "Stanford University",
			Description: "Engineering excellence",
			Logo:        "stanford_logo.png",
			BlurHash:    "blurhash123",
			StartDate:   time.Date(2016, 9, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC),
		},
		SchoolPeriods: []domain.SchoolPeriod{},
		Projects:      nil,
		Level:         domain.College,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	validResp, _ := json.Marshal(sampleEducation)

	type Given struct {
		method   string
		id       string
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
				method: http.MethodGet,
				id:     fixedID,
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Get(mock.Anything, fixedID).
						Return(sampleEducation, nil)
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
				id:       fixedID,
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusMethodNotAllowed,
				body: "Method not allowed: only GET is supported\n",
			},
		},
		"education not found": {
			given: Given{
				method: http.MethodGet,
				id:     "nonexistent",
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Get(mock.Anything, "nonexistent").
						Return(nil, nil)
				},
			},
			expected: Expected{
				code: http.StatusNotFound,
				body: "Education not found\n",
			},
		},
		"repository error": {
			given: Given{
				method: http.MethodGet,
				id:     fixedID,
				mockRepo: func(m *mockRepo.MockEducationRepository) {
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
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					// Testing Get() directly with empty id (bypasses ServeHTTP validation)
					m.EXPECT().
						Get(mock.Anything, "").
						Return(nil, nil)
				},
			},
			expected: Expected{
				code: http.StatusNotFound,
				body: "Education not found\n",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newEducationHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockEducationRepo)
			}

			req := httptest.NewRequest(tt.given.method, "/education/"+tt.given.id, nil)
			w := httptest.NewRecorder()

			f.educationHandler.Get(w, req, tt.given.id)

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

func TestEducationServiceHandler_Update(t *testing.T) {
	fixedID := "edu-123"

	validSchool := domain.SchoolPeriod{
		Name:        "Harvard University",
		Description: "Top-tier education",
		Logo:        "logo.png",
		BlurHash:    "hash123",
		StartDate:   time.Date(2015, 9, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
	}

	existingEducation := &domain.Education{
		Id:            fixedID,
		MainSchool:    validSchool,
		SchoolPeriods: []domain.SchoolPeriod{validSchool},
		Projects:      nil,
		Level:         domain.College,
		CreatedAt:     time.Now().Add(-time.Hour * 24),
		UpdatedAt:     time.Now(),
	}

	validBody, _ := json.Marshal(existingEducation)
	validResp, _ := json.Marshal(existingEducation)

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
				method: http.MethodPut,
				body:   string(validBody),
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Update(mock.Anything, mock.MatchedBy(func(edu *domain.Education) bool {
							return edu.Id == fixedID && edu.MainSchool.Name == "Harvard University"
						})).
						Return(existingEducation, nil)
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
					bad := *existingEducation
					bad.MainSchool.Name = "" // fails ValidatePayload()
					b, _ := json.Marshal(bad)
					return string(b)
				}(),
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid education payload: main school name missing\n",
			},
		},
		"repository error": {
			given: Given{
				method: http.MethodPut,
				body:   string(validBody),
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Update(mock.Anything, mock.AnythingOfType("*domain.Education")).
						Return(nil, errors.New("database failure"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to update education: database failure\n",
			},
		},
		"education not found": {
			given: Given{
				method: http.MethodPut,
				body:   string(validBody),
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Update(mock.Anything, mock.AnythingOfType("*domain.Education")).
						Return(nil, nil)
				},
			},
			expected: Expected{
				code: http.StatusNotFound,
				body: "Education not found\n",
			},
		},
		"large payload": {
			given: Given{
				method: http.MethodPut,
				body: func() string {
					large := *existingEducation
					large.MainSchool.Description = strings.Repeat("A", 10_000)
					b, _ := json.Marshal(large)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Update(mock.Anything, mock.AnythingOfType("*domain.Education")).
						Return(existingEducation, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: string(validResp),
			},
		},
		"unicode fields": {
			given: Given{
				method: http.MethodPut,
				body: func() string {
					unicode := *existingEducation
					unicode.MainSchool.Name = "東京大学 🏫"
					unicode.MainSchool.Description = "研究 excellence"
					b, _ := json.Marshal(unicode)
					return string(b)
				}(),
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Update(mock.Anything, mock.MatchedBy(func(edu *domain.Education) bool {
							return strings.Contains(edu.MainSchool.Name, "東京")
						})).
						Return(existingEducation, nil)
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
			f := newEducationHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockEducationRepo)
			}

			req := httptest.NewRequest(tt.given.method, "/education", strings.NewReader(tt.given.body))
			w := httptest.NewRecorder()

			f.educationHandler.Update(w, req)

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
