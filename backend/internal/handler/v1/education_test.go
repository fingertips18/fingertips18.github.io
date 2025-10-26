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

type educationHandlerTestFixture struct {
	t                 *testing.T
	mockEducationRepo *mockRepo.MockEducationRepository
	mockProjectRepo   *mockRepo.MockProjectRepository
	educationHandler  EducationHandler
}

func newEducationHandlerTestFixture(t *testing.T) *educationHandlerTestFixture {
	mockEducationRepo := new(mockRepo.MockEducationRepository)
	mockProjectRepo := new(mockRepo.MockProjectRepository)

	educationHandler := NewEducationServiceHandler(
		EducationServiceConfig{
			educationRepo: mockEducationRepo,
			ProjectRepo:   mockProjectRepo,
		},
	)

	return &educationHandlerTestFixture{
		t:                 t,
		mockEducationRepo: mockEducationRepo,
		mockProjectRepo:   mockProjectRepo,
		educationHandler:  educationHandler,
	}
}

func TestEducationServiceHandler_Create(t *testing.T) {
	fixedID := "edu-123"

	validSchool := SchoolPeriodDTO{
		Name:        "Harvard University",
		Description: "Top-tier education",
		Logo:        "logo.png",
		BlurHash:    "hash123",
		StartDate:   time.Date(2015, 9, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
	}

	validCreateReq := CreateEducationRequest{
		MainSchool:    validSchool,
		SchoolPeriods: []SchoolPeriodDTO{validSchool},
		Level:         "college",
	}
	validBody, _ := json.Marshal(validCreateReq)
	validResp, _ := json.Marshal(IDResponse{Id: fixedID})

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
					secondSchool := SchoolPeriodDTO{
						Name:        "Massachusetts Institute of Technology",
						Description: "Exchange program in Computer Science",
						Logo:        "mit_logo.png",
						BlurHash:    "mitHash456",
						StartDate:   time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
						EndDate:     time.Date(2018, 6, 1, 0, 0, 0, 0, time.UTC),
					}
					multi.SchoolPeriods = []SchoolPeriodDTO{
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

func TestEducationServiceHandler_Create_Routing(t *testing.T) {
	fixedID := "edu-123"

	validSchool := SchoolPeriodDTO{
		Name:        "Harvard University",
		Description: "Top-tier education",
		Logo:        "logo.png",
		BlurHash:    "hash123",
		StartDate:   time.Date(2015, 9, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
	}

	validCreateReq := EducationDTO{
		MainSchool:    validSchool,
		SchoolPeriods: []SchoolPeriodDTO{validSchool},
		Level:         "college",
	}
	validBody, _ := json.Marshal(validCreateReq)
	expectedResp, _ := json.Marshal(IDResponse{Id: fixedID})

	f := newEducationHandlerTestFixture(t)

	// Setup mock expectation
	f.mockEducationRepo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(edu *domain.Education) bool {
			return edu.Level == domain.College &&
				edu.MainSchool.Name == "Harvard University" &&
				len(edu.SchoolPeriods) == 1
		})).
		Return(fixedID, nil)

	// Create request and recorder
	req := httptest.NewRequest(http.MethodPost, "/education", strings.NewReader(string(validBody)))
	w := httptest.NewRecorder()

	// Cast handler to http.Handler and call ServeHTTP
	handler, ok := f.educationHandler.(http.Handler)
	assert.True(t, ok, "educationHandler should implement http.Handler")
	handler.ServeHTTP(w, req)

	// Verify response
	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.JSONEq(t, string(expectedResp), string(body))

	f.mockEducationRepo.AssertExpectations(t)
}

func TestEducationServiceHandler_Get(t *testing.T) {
	fixedID := "edu-123"
	fixedTime := time.Date(2025, 10, 25, 21, 56, 4, 0, time.UTC)

	sampleEducation := &EducationDTO{
		Id: fixedID,
		MainSchool: SchoolPeriodDTO{
			Name:        "Harvard University",
			Description: "Top-tier education",
			Logo:        "logo.png",
			BlurHash:    "hash123",
			StartDate:   time.Date(2015, 9, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
		},
		SchoolPeriods: []SchoolPeriodDTO{},
		Projects:      []ProjectDTO{},
		Level:         "college",
		CreatedAt:     fixedTime,
		UpdatedAt:     fixedTime,
	}

	validResp, _ := json.Marshal(sampleEducation)

	type Given struct {
		method       string
		id           string
		mockEducRepo func(m *mockRepo.MockEducationRepository)
		mockProjRepo func(m *mockRepo.MockProjectRepository)
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
				mockEducRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Get(mock.Anything, fixedID).
						Return(&domain.Education{
							Id: fixedID,
							MainSchool: domain.SchoolPeriod{
								Name:        "Harvard University",
								Description: "Top-tier education",
								Logo:        "logo.png",
								BlurHash:    "hash123",
								StartDate:   time.Date(2015, 9, 1, 0, 0, 0, 0, time.UTC),
								EndDate:     time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
							},
							SchoolPeriods: []domain.SchoolPeriod{},
							Level:         domain.College,
							CreatedAt:     fixedTime,
							UpdatedAt:     fixedTime,
						}, nil)
				},
				mockProjRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						ListByEducationID(mock.Anything, fixedID).
						Return([]domain.Project{}, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: string(validResp),
			},
		},
		"method not allowed": {
			given: Given{
				method:       http.MethodPost,
				id:           fixedID,
				mockEducRepo: nil,
				mockProjRepo: nil,
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
				mockEducRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Get(mock.Anything, "nonexistent").
						Return(nil, nil)
				},
				mockProjRepo: nil,
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
				mockEducRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Get(mock.Anything, fixedID).
						Return(nil, errors.New("database failure"))
				},
				mockProjRepo: nil,
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
				mockEducRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Get(mock.Anything, "").
						Return(nil, nil)
				},
				mockProjRepo: nil,
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

			if tt.given.mockEducRepo != nil {
				tt.given.mockEducRepo(f.mockEducationRepo)
			}
			if tt.given.mockProjRepo != nil {
				tt.given.mockProjRepo(f.mockProjectRepo)
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
			f.mockProjectRepo.AssertExpectations(t)
		})
	}
}

func TestEducationServiceHandler_Get_Routing(t *testing.T) {
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
		Level:         domain.College,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	dto := EducationDTO{
		Id: sampleEducation.Id,
		MainSchool: SchoolPeriodDTO{
			Name:        sampleEducation.MainSchool.Name,
			Description: sampleEducation.MainSchool.Description,
			Logo:        sampleEducation.MainSchool.Logo,
			BlurHash:    sampleEducation.MainSchool.BlurHash,
			Link:        sampleEducation.MainSchool.Link,
			Honor:       sampleEducation.MainSchool.Honor,
			StartDate:   sampleEducation.MainSchool.StartDate,
			EndDate:     sampleEducation.MainSchool.EndDate,
		},
		SchoolPeriods: func() []SchoolPeriodDTO {
			periods := make([]SchoolPeriodDTO, len(sampleEducation.SchoolPeriods))
			for i, p := range sampleEducation.SchoolPeriods {
				periods[i] = SchoolPeriodDTO{
					Name:        p.Name,
					Description: p.Description,
					Logo:        p.Logo,
					BlurHash:    p.BlurHash,
					Link:        p.Link,
					Honor:       p.Honor,
					StartDate:   p.StartDate,
					EndDate:     p.EndDate,
				}
			}
			return periods
		}(),
		Level:     string(sampleEducation.Level),
		CreatedAt: sampleEducation.CreatedAt,
		UpdatedAt: sampleEducation.UpdatedAt,
	}

	expectedResp, _ := json.Marshal(dto)

	f := newEducationHandlerTestFixture(t)

	// Setup mock expectation
	f.mockEducationRepo.EXPECT().
		Get(mock.Anything, fixedID).
		Return(sampleEducation, nil)

	// Setup project mock
	f.mockProjectRepo.EXPECT().
		ListByEducationID(mock.Anything, fixedID).
		Return([]domain.Project{}, nil)

	// Create request and recorder
	req := httptest.NewRequest(http.MethodGet, "/education/"+fixedID, nil)
	w := httptest.NewRecorder()

	// Cast handler to http.Handler and call ServeHTTP
	handler, ok := f.educationHandler.(http.Handler)
	assert.True(t, ok, "educationHandler should implement http.Handler")
	handler.ServeHTTP(w, req)

	// Verify response
	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, string(expectedResp), string(body))

	f.mockEducationRepo.AssertExpectations(t)
	f.mockProjectRepo.AssertExpectations(t)
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
		Level:         domain.College,
		CreatedAt:     time.Now().Add(-time.Hour * 24),
		UpdatedAt:     time.Now(),
	}

	dto := EducationDTO{
		Id: existingEducation.Id,
		MainSchool: SchoolPeriodDTO{
			Link:        existingEducation.MainSchool.Link,
			Name:        existingEducation.MainSchool.Name,
			Description: existingEducation.MainSchool.Description,
			Logo:        existingEducation.MainSchool.Logo,
			BlurHash:    existingEducation.MainSchool.BlurHash,
			Honor:       existingEducation.MainSchool.Honor,
			StartDate:   existingEducation.MainSchool.StartDate,
			EndDate:     existingEducation.MainSchool.EndDate,
		},
		SchoolPeriods: func() []SchoolPeriodDTO {
			periods := make([]SchoolPeriodDTO, len(existingEducation.SchoolPeriods))
			for i, p := range existingEducation.SchoolPeriods {
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
		Level:     string(existingEducation.Level),
		CreatedAt: existingEducation.CreatedAt,
		UpdatedAt: existingEducation.UpdatedAt,
	}

	validBody, _ := json.Marshal(dto)
	validResp, _ := json.Marshal(dto)

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
					bad.MainSchool = domain.SchoolPeriod{} // fails ValidatePayload()
					b, _ := json.Marshal(bad)
					return string(b)
				}(),
				mockRepo: nil,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid education payload: main school missing\n",
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
					large := dto
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
					unicode := dto
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

func TestEducationServiceHandler_Update_Routing(t *testing.T) {
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
		Level:         domain.College,
		CreatedAt:     time.Now().Add(-time.Hour * 24),
		UpdatedAt:     time.Now(),
	}

	dto := EducationDTO{
		Id: existingEducation.Id,
		MainSchool: SchoolPeriodDTO{
			Name:        existingEducation.MainSchool.Name,
			Description: existingEducation.MainSchool.Description,
			Logo:        existingEducation.MainSchool.Logo,
			BlurHash:    existingEducation.MainSchool.BlurHash,
			Link:        existingEducation.MainSchool.Link,
			Honor:       existingEducation.MainSchool.Honor,
			StartDate:   existingEducation.MainSchool.StartDate,
			EndDate:     existingEducation.MainSchool.EndDate,
		},
		SchoolPeriods: func() []SchoolPeriodDTO {
			periods := make([]SchoolPeriodDTO, len(existingEducation.SchoolPeriods))
			for i, p := range existingEducation.SchoolPeriods {
				periods[i] = SchoolPeriodDTO{
					Name:        p.Name,
					Description: p.Description,
					Logo:        p.Logo,
					BlurHash:    p.BlurHash,
					Link:        p.Link,
					Honor:       p.Honor,
					StartDate:   p.StartDate,
					EndDate:     p.EndDate,
				}
			}
			return periods
		}(),
		Level:     string(existingEducation.Level),
		CreatedAt: existingEducation.CreatedAt,
		UpdatedAt: existingEducation.UpdatedAt,
	}

	validBody, _ := json.Marshal(dto)
	expectedResp, _ := json.Marshal(dto)

	f := newEducationHandlerTestFixture(t)

	// Setup mock expectation
	f.mockEducationRepo.EXPECT().
		Update(mock.Anything, mock.MatchedBy(func(edu *domain.Education) bool {
			return edu.Id == fixedID && edu.MainSchool.Name == "Harvard University"
		})).
		Return(existingEducation, nil)

	// Create request and recorder
	req := httptest.NewRequest(http.MethodPut, "/education", strings.NewReader(string(validBody)))
	w := httptest.NewRecorder()

	// Cast handler to http.Handler and call ServeHTTP
	handler, ok := f.educationHandler.(http.Handler)
	assert.True(t, ok, "educationHandler should implement http.Handler")
	handler.ServeHTTP(w, req)

	// Verify response
	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, string(expectedResp), string(body))

	f.mockEducationRepo.AssertExpectations(t)
}

func TestEducationServiceHandler_Delete(t *testing.T) {
	fixedID := "edu-123"

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
				method: http.MethodDelete,
				id:     fixedID,
				mockRepo: func(m *mockRepo.MockEducationRepository) {
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
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					// no call expected
				},
			},
			expected: Expected{
				code: http.StatusMethodNotAllowed,
				body: "Method not allowed: only DELETE is supported\n",
			},
		},
		"education not found (pgx.ErrNoRows)": {
			given: Given{
				method: http.MethodDelete,
				id:     "missing-id",
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Delete(mock.Anything, "missing-id").
						Return(pgx.ErrNoRows)
				},
			},
			expected: Expected{
				code: http.StatusNotFound,
				body: "Education not found\n",
			},
		},
		"repository error": {
			given: Given{
				method: http.MethodDelete,
				id:     fixedID,
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Delete(mock.Anything, fixedID).
						Return(errors.New("database failure"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to delete education: database failure\n",
			},
		},
		"empty id": {
			given: Given{
				method: http.MethodDelete,
				id:     "",
				mockRepo: func(m *mockRepo.MockEducationRepository) {
					m.EXPECT().
						Delete(mock.Anything, "").
						Return(pgx.ErrNoRows)
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

			f.educationHandler.Delete(w, req, tt.given.id)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)
			assert.Equal(t, tt.expected.body, string(body))

			f.mockEducationRepo.AssertExpectations(t)
		})
	}
}

func TestEducationServiceHandler_Delete_Routing(t *testing.T) {
	fixedID := "edu-123"

	f := newEducationHandlerTestFixture(t)

	// Setup mock expectation
	f.mockEducationRepo.EXPECT().
		Delete(mock.Anything, fixedID).
		Return(nil)

	// Create request and recorder
	req := httptest.NewRequest(http.MethodDelete, "/education/"+fixedID, nil)
	w := httptest.NewRecorder()

	// Cast handler to http.Handler and call ServeHTTP
	handler, ok := f.educationHandler.(http.Handler)
	assert.True(t, ok, "educationHandler should implement http.Handler")
	handler.ServeHTTP(w, req)

	// Verify response
	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	assert.Equal(t, "", string(body))

	f.mockEducationRepo.AssertExpectations(t)
}

func TestEducationServiceHandler_List(t *testing.T) {
	sampleEducation := domain.Education{
		Id: "edu-123",
		MainSchool: domain.SchoolPeriod{
			Name:        "MIT",
			Description: "Computer Science",
			Logo:        "mit.png",
			BlurHash:    "hash123",
			StartDate:   time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2018, 6, 1, 0, 0, 0, 0, time.UTC),
		},
		SchoolPeriods: []domain.SchoolPeriod{},
		Level:         domain.College,
		CreatedAt:     time.Now().Add(-24 * time.Hour),
		UpdatedAt:     time.Now(),
	}

	listResp := []domain.Education{sampleEducation}

	educations := make([]EducationDTO, len(listResp))
	for i, e := range listResp {
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
			Projects:  []ProjectDTO{}, // Empty projects array
			Level:     string(e.Level),
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}
	}

	validJSON, _ := json.Marshal(educations)

	type Given struct {
		method       string
		query        string
		mockEducRepo func(m *mockRepo.MockEducationRepository)
		mockProjRepo func(m *mockRepo.MockProjectRepository)
	}
	type Expected struct {
		code int
		body string
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"success with default params": {
			given: Given{
				method: http.MethodGet,
				query:  "",
				mockEducRepo: func(m *mockRepo.MockEducationRepository) {
					expectedFilter := domain.EducationFilter{
						Page:          1,
						PageSize:      10,
						SortBy:        nil,
						SortAscending: false,
					}
					m.EXPECT().
						List(mock.Anything, expectedFilter).
						Return(listResp, nil)
				},
				mockProjRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						ListByEducationID(mock.Anything, "edu-123").
						Return([]domain.Project{}, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: string(validJSON),
			},
		},
		"method not allowed": {
			given: Given{
				method:       http.MethodPost,
				query:        "",
				mockEducRepo: nil,
				mockProjRepo: nil,
			},
			expected: Expected{
				code: http.StatusMethodNotAllowed,
				body: "Method not allowed: only GET is supported\n",
			},
		},
		"invalid sort_by": {
			given: Given{
				method:       http.MethodGet,
				query:        "?sort_by=!!invalid",
				mockEducRepo: nil,
				mockProjRepo: nil,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "invalid sort by\n",
			},
		},
		"repository error": {
			given: Given{
				method: http.MethodGet,
				query:  "?page=1&page_size=10&sort_by=created_at",
				mockEducRepo: func(m *mockRepo.MockEducationRepository) {
					sortBy := domain.CreatedAt
					expectedFilter := domain.EducationFilter{
						Page:          1,
						PageSize:      10,
						SortBy:        &sortBy,
						SortAscending: false,
					}
					m.EXPECT().
						List(mock.Anything, expectedFilter).
						Return(nil, errors.New("database failure"))
				},
				mockProjRepo: nil, // No project repo call when education repo fails
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to list educations: database failure\n",
			},
		},
		"success with custom filters": {
			given: Given{
				method: http.MethodGet,
				query:  "?page=2&page_size=5&sort_by=updated_at&sort_ascending=true",
				mockEducRepo: func(m *mockRepo.MockEducationRepository) {
					sortBy := domain.UpdatedAt
					expectedFilter := domain.EducationFilter{
						Page:          2,
						PageSize:      5,
						SortBy:        &sortBy,
						SortAscending: true,
					}
					m.EXPECT().
						List(mock.Anything, expectedFilter).
						Return(listResp, nil)
				},
				mockProjRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						ListByEducationID(mock.Anything, "edu-123").
						Return([]domain.Project{}, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: string(validJSON),
			},
		},
		"empty list response": {
			given: Given{
				method: http.MethodGet,
				query:  "?sort_by=created_at",
				mockEducRepo: func(m *mockRepo.MockEducationRepository) {
					sortBy := domain.CreatedAt
					expectedFilter := domain.EducationFilter{
						Page:          1,
						PageSize:      10,
						SortBy:        &sortBy,
						SortAscending: false,
					}
					m.EXPECT().
						List(mock.Anything, expectedFilter).
						Return([]domain.Education{}, nil)
				},
				mockProjRepo: nil, // No project repo call when list is empty
			},
			expected: Expected{
				code: http.StatusOK,
				body: "[]",
			},
		},
		"page zero defaults to one": {
			given: Given{
				method: http.MethodGet,
				query:  "?page=0",
				mockEducRepo: func(m *mockRepo.MockEducationRepository) {
					expectedFilter := domain.EducationFilter{
						Page:          1,
						PageSize:      10,
						SortBy:        nil,
						SortAscending: false,
					}
					m.EXPECT().
						List(mock.Anything, expectedFilter).
						Return(listResp, nil)
				},
				mockProjRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						ListByEducationID(mock.Anything, "edu-123").
						Return([]domain.Project{}, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: string(validJSON),
			},
		},
		"negative page size": {
			given: Given{
				method: http.MethodGet,
				query:  "?page_size=-1",
				mockEducRepo: func(m *mockRepo.MockEducationRepository) {
					expectedFilter := domain.EducationFilter{
						Page:          1,
						PageSize:      10,
						SortBy:        nil,
						SortAscending: false,
					}
					m.EXPECT().
						List(mock.Anything, expectedFilter).
						Return(listResp, nil)
				},
				mockProjRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						ListByEducationID(mock.Anything, "edu-123").
						Return([]domain.Project{}, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: string(validJSON),
			},
		},
		"page size exceeds maximum": {
			given: Given{
				method: http.MethodGet,
				query:  "?page_size=1000",
				mockEducRepo: func(m *mockRepo.MockEducationRepository) {
					expectedFilter := domain.EducationFilter{
						Page:          1,
						PageSize:      100,
						SortBy:        nil,
						SortAscending: false,
					}
					m.EXPECT().
						List(mock.Anything, expectedFilter).
						Return(listResp, nil)
				},
				mockProjRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						ListByEducationID(mock.Anything, "edu-123").
						Return([]domain.Project{}, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: string(validJSON),
			},
		},
		"invalid page parameter": {
			given: Given{
				method: http.MethodGet,
				query:  "?page=invalid",
				mockEducRepo: func(m *mockRepo.MockEducationRepository) {
					expectedFilter := domain.EducationFilter{
						Page:          1,
						PageSize:      10,
						SortBy:        nil,
						SortAscending: false,
					}
					m.EXPECT().
						List(mock.Anything, expectedFilter).
						Return(listResp, nil)
				},
				mockProjRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						ListByEducationID(mock.Anything, "edu-123").
						Return([]domain.Project{}, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: string(validJSON),
			},
		},
		"invalid page size parameter": {
			given: Given{
				method: http.MethodGet,
				query:  "?page_size=invalid",
				mockEducRepo: func(m *mockRepo.MockEducationRepository) {
					expectedFilter := domain.EducationFilter{
						Page:          1,
						PageSize:      10,
						SortBy:        nil,
						SortAscending: false,
					}
					m.EXPECT().
						List(mock.Anything, expectedFilter).
						Return(listResp, nil)
				},
				mockProjRepo: func(m *mockRepo.MockProjectRepository) {
					m.EXPECT().
						ListByEducationID(mock.Anything, "edu-123").
						Return([]domain.Project{}, nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: string(validJSON),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newEducationHandlerTestFixture(t)
			if tt.given.mockEducRepo != nil {
				tt.given.mockEducRepo(f.mockEducationRepo)
			}
			if tt.given.mockProjRepo != nil {
				tt.given.mockProjRepo(f.mockProjectRepo)
			}

			req := httptest.NewRequest(tt.given.method, "/educations"+tt.given.query, nil)
			w := httptest.NewRecorder()

			f.educationHandler.List(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)

			if strings.HasPrefix(tt.expected.body, "[") {
				assert.JSONEq(t, tt.expected.body, string(body))
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			f.mockEducationRepo.AssertExpectations(t)
			f.mockProjectRepo.AssertExpectations(t)
		})
	}
}

func TestEducationServiceHandler_List_Routing(t *testing.T) {
	sampleEducation := domain.Education{
		Id: "edu-123",
		MainSchool: domain.SchoolPeriod{
			Name:        "MIT",
			Description: "Computer Science",
			Logo:        "mit.png",
			BlurHash:    "hash123",
			StartDate:   time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2018, 6, 1, 0, 0, 0, 0, time.UTC),
		},
		SchoolPeriods: []domain.SchoolPeriod{},
		Level:         domain.College,
		CreatedAt:     time.Now().Add(-24 * time.Hour),
		UpdatedAt:     time.Now(),
	}

	listResp := []domain.Education{sampleEducation}

	educations := make([]EducationDTO, len(listResp))
	for i, e := range listResp {
		educations[i] = EducationDTO{
			Id: e.Id,
			MainSchool: SchoolPeriodDTO{
				Name:        e.MainSchool.Name,
				Description: e.MainSchool.Description,
				Logo:        e.MainSchool.Logo,
				BlurHash:    e.MainSchool.BlurHash,
				StartDate:   e.MainSchool.StartDate,
				EndDate:     e.MainSchool.EndDate,
				Link:        e.MainSchool.Link,
				Honor:       e.MainSchool.Honor,
			},
			SchoolPeriods: func() []SchoolPeriodDTO {
				periods := make([]SchoolPeriodDTO, len(e.SchoolPeriods))
				for j, p := range e.SchoolPeriods {
					periods[j] = SchoolPeriodDTO{
						Name:        p.Name,
						Description: p.Description,
						Logo:        p.Logo,
						BlurHash:    p.BlurHash,
						StartDate:   p.StartDate,
						EndDate:     p.EndDate,
						Link:        p.Link,
						Honor:       p.Honor,
					}
				}
				return periods
			}(),
			Projects:  []ProjectDTO{}, // Add this field
			Level:     string(e.Level),
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}
	}

	validJSON, _ := json.Marshal(educations)

	f := newEducationHandlerTestFixture(t)

	// Setup mock expectation for education repository
	expectedFilter := domain.EducationFilter{
		Page:          1,
		PageSize:      10,
		SortBy:        nil,
		SortAscending: false,
	}
	f.mockEducationRepo.EXPECT().
		List(mock.Anything, expectedFilter).
		Return(listResp, nil)

	// Setup mock expectation for project repository
	f.mockProjectRepo.EXPECT().
		ListByEducationID(mock.Anything, "edu-123").
		Return([]domain.Project{}, nil)

	// Create request and recorder
	req := httptest.NewRequest(http.MethodGet, "/educations", nil)
	w := httptest.NewRecorder()

	// Cast handler to http.Handler and call ServeHTTP
	handler, ok := f.educationHandler.(http.Handler)
	assert.True(t, ok, "educationHandler should implement http.Handler")
	handler.ServeHTTP(w, req)

	// Verify response
	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, string(validJSON), string(body))

	f.mockEducationRepo.AssertExpectations(t)
	f.mockProjectRepo.AssertExpectations(t) // Add this assertion
}
