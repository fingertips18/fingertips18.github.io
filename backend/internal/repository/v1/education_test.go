package v1

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	database "github.com/fingertips18/fingertips18.github.io/backend/internal/database/mocks"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testEducationTable = "test-education"
)

type educationRepositoryTestFixture struct {
	t                   *testing.T
	databaseAPI         *database.MockDatabaseAPI
	educationRepository educationRepository
}

func newEducationRepositoryTestFixture(t *testing.T, timeProvider func() time.Time) *educationRepositoryTestFixture {
	mockDatabaseAPI := new(database.MockDatabaseAPI)
	educationRepository := &educationRepository{
		databaseAPI:    mockDatabaseAPI,
		timeProvider:   timeProvider,
		educationTable: testEducationTable,
	}

	return &educationRepositoryTestFixture{
		t:                   t,
		databaseAPI:         mockDatabaseAPI,
		educationRepository: *educationRepository,
	}
}

func TestEducationRepository_Create(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	scanErr := errors.New("scan error")

	validEducation := domain.Education{
		MainSchool: domain.SchoolPeriod{
			Link:        "http://example.com",
			Name:        "test-name",
			Description: "test-description",
			Logo:        "test-logo",
			BlurHash:    "test-blurhash",
			Honor:       "test-honor",
			StartDate:   time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC),
			EndDate:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		Projects: []domain.Project{
			{
				Preview:     "test-preview",
				BlurHash:    "test-blurhash",
				Title:       "test-title",
				SubTitle:    "test-subtitle",
				Description: "test-description",
				Stack:       []string{"stack1"},
				Type:        domain.Web,
				Link:        "http://example.com",
				CreatedAt:   fixedTime,
				UpdatedAt:   fixedTime,
			},
		},
		Level:     domain.College,
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}

	type Given struct {
		education    domain.Education
		mockQueryRow func(m *database.MockDatabaseAPI)
	}

	type Expected struct {
		err error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful create education": {
			given: Given{
				education: validEducation,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).Return(&fakeRow{id: "123-abc"})
				},
			},
			expected: Expected{
				err: nil,
			},
		},
		"Database scan fails": {
			given: Given{
				education: validEducation,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(mock.Anything, mock.Anything, mock.Anything).
						Return(&fakeRow{scanErr: scanErr})
				},
			},
			expected: Expected{
				err: fmt.Errorf("failed to create education: %w", scanErr),
			},
		},
		"Database returns empty ID": {
			given: Given{
				education: validEducation,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, mock.Anything).
						Return(&fakeRow{id: ""})
				},
			},
			expected: Expected{
				err: errors.New("invalid education returned: ID missing"),
			},
		},
		"Missing main school fails": {
			given: Given{
				education: domain.Education{
					MainSchool:    domain.SchoolPeriod{},
					SchoolPeriods: validEducation.SchoolPeriods,
					Projects:      validEducation.Projects,
					Level:         validEducation.Level,
					CreatedAt:     validEducation.CreatedAt,
					UpdatedAt:     validEducation.UpdatedAt,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: main school missing"),
			},
		},
		"Missing main school name fails": {
			given: Given{
				education: domain.Education{
					MainSchool: domain.SchoolPeriod{
						Name:        "",
						Description: validEducation.MainSchool.Description,
						Logo:        validEducation.MainSchool.Logo,
						BlurHash:    validEducation.MainSchool.BlurHash,
						StartDate:   validEducation.MainSchool.StartDate,
						EndDate:     validEducation.MainSchool.EndDate,
					},
					SchoolPeriods: validEducation.SchoolPeriods,
					Projects:      validEducation.Projects,
					Level:         validEducation.Level,
					CreatedAt:     validEducation.CreatedAt,
					UpdatedAt:     validEducation.UpdatedAt,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: main school name missing"),
			},
		},
		"Missing main school description fails": {
			given: Given{
				education: domain.Education{
					MainSchool: domain.SchoolPeriod{
						Name:        validEducation.MainSchool.Name,
						Description: "",
						Logo:        validEducation.MainSchool.Logo,
						BlurHash:    validEducation.MainSchool.BlurHash,
						StartDate:   validEducation.MainSchool.StartDate,
						EndDate:     validEducation.MainSchool.EndDate,
					},
					SchoolPeriods: validEducation.SchoolPeriods,
					Projects:      validEducation.Projects,
					Level:         validEducation.Level,
					CreatedAt:     validEducation.CreatedAt,
					UpdatedAt:     validEducation.UpdatedAt,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: main school description missing"),
			},
		},
		"Missing main school logo fails": {
			given: Given{
				education: domain.Education{
					MainSchool: domain.SchoolPeriod{
						Name:        validEducation.MainSchool.Name,
						Description: validEducation.MainSchool.Description,
						Logo:        "",
						BlurHash:    validEducation.MainSchool.BlurHash,
						StartDate:   validEducation.MainSchool.StartDate,
						EndDate:     validEducation.MainSchool.EndDate,
					},
					SchoolPeriods: validEducation.SchoolPeriods,
					Projects:      validEducation.Projects,
					Level:         validEducation.Level,
					CreatedAt:     validEducation.CreatedAt,
					UpdatedAt:     validEducation.UpdatedAt,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: main school logo missing"),
			},
		},
		"Missing main school blurhash fails": {
			given: Given{
				education: domain.Education{
					MainSchool: domain.SchoolPeriod{
						Name:        validEducation.MainSchool.Name,
						Description: validEducation.MainSchool.Description,
						Logo:        validEducation.MainSchool.Logo,
						BlurHash:    "",
						StartDate:   validEducation.MainSchool.StartDate,
						EndDate:     validEducation.MainSchool.EndDate,
					},
					SchoolPeriods: validEducation.SchoolPeriods,
					Projects:      validEducation.Projects,
					Level:         validEducation.Level,
					CreatedAt:     validEducation.CreatedAt,
					UpdatedAt:     validEducation.UpdatedAt,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: main school blurHash missing"),
			},
		},
		"Missing main school start date fails": {
			given: Given{
				education: domain.Education{
					MainSchool: domain.SchoolPeriod{
						Name:        validEducation.MainSchool.Name,
						Description: validEducation.MainSchool.Description,
						Logo:        validEducation.MainSchool.Logo,
						BlurHash:    validEducation.MainSchool.BlurHash,
						StartDate:   time.Time{},
						EndDate:     validEducation.MainSchool.EndDate,
					},
					SchoolPeriods: validEducation.SchoolPeriods,
					Projects:      validEducation.Projects,
					Level:         validEducation.Level,
					CreatedAt:     validEducation.CreatedAt,
					UpdatedAt:     validEducation.UpdatedAt,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: main school start date missing"),
			},
		},
		"Missing main school end date fails": {
			given: Given{
				education: domain.Education{
					MainSchool: domain.SchoolPeriod{
						Name:        validEducation.MainSchool.Name,
						Description: validEducation.MainSchool.Description,
						Logo:        validEducation.MainSchool.Logo,
						BlurHash:    validEducation.MainSchool.BlurHash,
						StartDate:   validEducation.MainSchool.StartDate,
						EndDate:     time.Time{},
					},
					SchoolPeriods: validEducation.SchoolPeriods,
					Projects:      validEducation.Projects,
					Level:         validEducation.Level,
					CreatedAt:     validEducation.CreatedAt,
					UpdatedAt:     validEducation.UpdatedAt,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: main school end date missing"),
			},
		},
		"Missing school period fails": {
			given: Given{
				education: domain.Education{
					MainSchool: validEducation.MainSchool,
					SchoolPeriods: []domain.SchoolPeriod{
						{},
					},
					Projects:  validEducation.Projects,
					Level:     validEducation.Level,
					CreatedAt: validEducation.CreatedAt,
					UpdatedAt: validEducation.UpdatedAt,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: school period[0] is empty"),
			},
		},
		"Missing school period value fails": {
			given: Given{
				education: domain.Education{
					MainSchool: validEducation.MainSchool,
					SchoolPeriods: []domain.SchoolPeriod{
						{
							Name:        "",
							Description: validEducation.MainSchool.Description,
							Logo:        validEducation.MainSchool.Logo,
							BlurHash:    validEducation.MainSchool.BlurHash,
							StartDate:   validEducation.MainSchool.StartDate,
							EndDate:     validEducation.MainSchool.EndDate,
						},
					},
					Projects:  validEducation.Projects,
					Level:     validEducation.Level,
					CreatedAt: validEducation.CreatedAt,
					UpdatedAt: validEducation.UpdatedAt,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: school period[0] name missing"),
			},
		},
		"Missing projects value fails": {
			given: Given{
				education: domain.Education{
					MainSchool: validEducation.MainSchool,
					Projects: []domain.Project{
						{
							Preview:     "",
							BlurHash:    validEducation.Projects[0].BlurHash,
							Title:       validEducation.Projects[0].Title,
							SubTitle:    validEducation.Projects[0].SubTitle,
							Description: validEducation.Projects[0].Description,
							Stack:       validEducation.Projects[0].Stack,
							Type:        validEducation.Projects[0].Type,
							Link:        validEducation.Projects[0].Link,
						},
					},
					Level:     validEducation.Level,
					CreatedAt: validEducation.CreatedAt,
					UpdatedAt: validEducation.UpdatedAt,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: project[0] preview missing"),
			},
		},
		"Missing level fails": {
			given: Given{
				education: domain.Education{
					MainSchool: validEducation.MainSchool,
					Projects:   validEducation.Projects,
					Level:      "",
					CreatedAt:  validEducation.CreatedAt,
					UpdatedAt:  validEducation.UpdatedAt,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: level missing"),
			},
		},
		"Invalid level fails": {
			given: Given{
				education: domain.Education{
					MainSchool: validEducation.MainSchool,
					Projects:   validEducation.Projects,
					Level:      "invalid",
					CreatedAt:  validEducation.CreatedAt,
					UpdatedAt:  validEducation.UpdatedAt,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: level invalid = invalid"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newEducationRepositoryTestFixture(t, func() time.Time { return fixedTime })

			if test.given.mockQueryRow != nil {
				test.given.mockQueryRow(f.databaseAPI)
			}

			id, err := f.educationRepository.Create(context.Background(), test.given.education)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Empty(t, err)
				assert.NotEmpty(t, id)
				assert.Equal(t, fixedTime, test.given.education.CreatedAt)
				assert.Equal(t, fixedTime, test.given.education.UpdatedAt)
			}

			f.databaseAPI.AssertExpectations(t)
		})

	}
}
