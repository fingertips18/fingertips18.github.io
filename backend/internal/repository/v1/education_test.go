package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	database "github.com/fingertips18/fingertips18.github.io/backend/internal/database/mocks"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testEducationTable = "test-education"
)

type educationCreateFakeRow struct {
	id      string
	scanErr error
}

func (r *educationCreateFakeRow) Scan(dest ...interface{}) error {
	if r.scanErr != nil {
		return r.scanErr
	}
	if len(dest) != 1 {
		return fmt.Errorf("expected 1 scan destination, got %d", len(dest))
	}
	if ptr, ok := dest[0].(*string); ok {
		*ptr = r.id
	}
	return nil
}

type educationGetFakeRow struct {
	id                string
	mainSchoolJSON    []byte
	schoolPeriodsJSON []byte
	projectsJSON      []byte
	level             domain.EducationLevel
	createdAt         time.Time
	updatedAt         time.Time
	scanErr           error
}

func (r *educationGetFakeRow) Scan(dest ...interface{}) error {
	if r.scanErr != nil {
		return r.scanErr
	}

	if len(dest) != 7 {
		return fmt.Errorf("expected 7 scan destinations, got %d", len(dest))
	}

	if ptr, ok := dest[0].(*string); ok {
		*ptr = r.id
	}
	if ptr, ok := dest[1].(*[]byte); ok {
		*ptr = r.mainSchoolJSON
	}
	if ptr, ok := dest[2].(*[]byte); ok {
		*ptr = r.schoolPeriodsJSON
	}
	if ptr, ok := dest[3].(*[]byte); ok {
		*ptr = r.projectsJSON
	}
	if ptr, ok := dest[4].(*domain.EducationLevel); ok {
		*ptr = r.level
	}
	if ptr, ok := dest[5].(*time.Time); ok {
		*ptr = r.createdAt
	}
	if ptr, ok := dest[6].(*time.Time); ok {
		*ptr = r.updatedAt
	}

	return nil
}

type educationFakeCommandTag string

func (f educationFakeCommandTag) RowsAffected() int64 {
	if f == "DELETE 1" {
		return 1
	}
	return 0
}

type educationFakeRow struct {
	education domain.Education
	scanErr   error
}
type educationFakeRows struct {
	rows   []*educationFakeRow
	index  int
	rowErr error
}

func (r *educationFakeRows) Next() bool {
	if r.index >= len(r.rows) {
		return false
	}
	r.index++
	return true
}

func (r *educationFakeRows) Scan(dest ...interface{}) error {
	row := r.rows[r.index-1]
	if row.scanErr != nil {
		return row.scanErr
	}
	ed := row.education

	*(dest[0].(*string)) = ed.Id
	*(dest[1].(*domain.SchoolPeriod)) = ed.MainSchool
	*(dest[2].(*[]domain.SchoolPeriod)) = ed.SchoolPeriods
	*(dest[3].(*[]domain.Project)) = ed.Projects
	*(dest[4].(*domain.EducationLevel)) = ed.Level
	*(dest[5].(*time.Time)) = ed.CreatedAt
	*(dest[6].(*time.Time)) = ed.UpdatedAt

	return nil
}

func (r *educationFakeRows) Close() {}

func (r *educationFakeRows) Err() error {
	return r.rowErr
}

type educationRepositoryTestFixture struct {
	t                   *testing.T
	databaseAPI         *database.MockDatabaseAPI
	educationRepository *educationRepository
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
		educationRepository: educationRepository,
	}
}

func TestEducationRepository_Create(t *testing.T) {
	fixedID := "123-abc"
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
			},
		},
		Level: domain.College,
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
					).Return(&educationCreateFakeRow{id: fixedID})
				},
			},
			expected: Expected{
				err: nil,
			},
		},
		"Valid multiple school periods": {
			given: Given{
				education: domain.Education{
					MainSchool: validEducation.MainSchool,
					SchoolPeriods: []domain.SchoolPeriod{
						validEducation.MainSchool,
						{
							Name:        "Another School",
							Description: "Desc",
							Logo:        "Logo2",
							BlurHash:    "BlurHash2",
							StartDate:   time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
							EndDate:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					Projects: validEducation.Projects,
					Level:    validEducation.Level,
				},
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(mock.Anything, mock.Anything, mock.Anything).Return(&educationCreateFakeRow{id: fixedID})
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
						Return(&educationCreateFakeRow{scanErr: scanErr})
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
						Return(&educationCreateFakeRow{id: ""})
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
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: main school end date missing"),
			},
		},
		"End date before start date fails": {
			given: Given{
				education: domain.Education{
					MainSchool: domain.SchoolPeriod{
						Name:        validEducation.MainSchool.Name,
						Description: validEducation.MainSchool.Description,
						Logo:        validEducation.MainSchool.Logo,
						BlurHash:    validEducation.MainSchool.BlurHash,
						StartDate:   time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
						EndDate:     time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
					},
					SchoolPeriods: validEducation.SchoolPeriods,
					Projects:      validEducation.Projects,
					Level:         validEducation.Level,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: main school end date must be after start date"),
			},
		},
		"Missing school period fails": {
			given: Given{
				education: domain.Education{
					MainSchool: validEducation.MainSchool,
					SchoolPeriods: []domain.SchoolPeriod{
						{},
					},
					Projects: validEducation.Projects,
					Level:    validEducation.Level,
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
					Projects: validEducation.Projects,
					Level:    validEducation.Level,
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
					Level: validEducation.Level,
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

			id, err := f.educationRepository.Create(context.Background(), &test.given.education)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, fixedID, id)
				assert.Equal(t, fixedTime, test.given.education.CreatedAt)
				assert.Equal(t, fixedTime, test.given.education.UpdatedAt)
			}

			f.databaseAPI.AssertExpectations(t)
		})

	}
}

func TestEducationRepository_Get(t *testing.T) {
	fixedID := "123-abc"
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	scanErr := errors.New("scan error")

	validMainSchool := domain.SchoolPeriod{
		Link:        "http://example.com",
		Name:        "test-name",
		Description: "test-description",
		Logo:        "test-logo",
		BlurHash:    "test-blurhash",
		Honor:       "test-honor",
		StartDate:   time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	validProjects := []domain.Project{
		{
			Preview:     "test-preview",
			BlurHash:    "test-blurhash",
			Title:       "test-title",
			SubTitle:    "test-subtitle",
			Description: "test-description",
			Stack:       []string{"stack1"},
			Type:        domain.Web,
			Link:        "http://example.com",
		},
	}

	mainSchoolJSON, _ := json.Marshal(validMainSchool)
	schoolPeriodsJSON, _ := json.Marshal([]domain.SchoolPeriod{validMainSchool})
	projectsJSON, _ := json.Marshal(validProjects)

	type Given struct {
		id           string
		mockQueryRow func(m *database.MockDatabaseAPI)
	}

	type Expected struct {
		education *domain.Education
		err       error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful get education": {
			given: Given{
				id: fixedID,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(mock.Anything, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
						return len(args) == 1 && args[0] == fixedID
					})).
						Return(&educationGetFakeRow{
							id:                fixedID,
							mainSchoolJSON:    mainSchoolJSON,
							schoolPeriodsJSON: schoolPeriodsJSON,
							projectsJSON:      projectsJSON,
							level:             domain.College,
							createdAt:         fixedTime,
							updatedAt:         fixedTime,
						})
				},
			},
			expected: Expected{
				education: &domain.Education{
					Id:            fixedID,
					MainSchool:    validMainSchool,
					SchoolPeriods: []domain.SchoolPeriod{validMainSchool},
					Projects:      validProjects,
					Level:         domain.College,
					CreatedAt:     fixedTime,
					UpdatedAt:     fixedTime,
				},
				err: nil,
			},
		},
		"Empty ID": {
			given: Given{
				id:           "",
				mockQueryRow: nil,
			},
			expected: Expected{
				education: nil,
				err:       errors.New("failed to get education: ID missing"),
			},
		},
		"Row not found": {
			given: Given{
				id: fixedID,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(mock.Anything, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
						return len(args) == 1 && args[0] == fixedID
					})).
						Return(&educationGetFakeRow{scanErr: pgx.ErrNoRows})
				},
			},
			expected: Expected{
				education: nil,
				err:       fmt.Errorf("education not found: %w", pgx.ErrNoRows),
			},
		},
		"Database scan error": {
			given: Given{
				id: fixedID,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(mock.Anything, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
						return len(args) == 1 && args[0] == fixedID
					})).
						Return(&educationGetFakeRow{scanErr: scanErr})
				},
			},
			expected: Expected{
				education: nil,
				err:       fmt.Errorf("failed to scan education: %w", scanErr),
			},
		},
		"Invalid main school JSON": {
			given: Given{
				id: fixedID,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(mock.Anything, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
						return len(args) == 1 && args[0] == fixedID
					})).
						Return(&educationGetFakeRow{
							id:             fixedID,
							mainSchoolJSON: []byte("{invalid-json}"),
							level:          domain.College,
							createdAt:      fixedTime,
							updatedAt:      fixedTime,
						})
				},
			},
			expected: Expected{
				education: nil,
				err:       fmt.Errorf("failed to unmarshal main school: %w", errors.New("invalid character 'i' looking for beginning of object key string")),
			},
		},
		"Validation error on level": {
			given: Given{
				id: fixedID,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(mock.Anything, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
						return len(args) == 1 && args[0] == fixedID
					})).
						Return(&educationGetFakeRow{
							id:             fixedID,
							mainSchoolJSON: mainSchoolJSON,
							level:          "invalid-level",
							createdAt:      fixedTime,
							updatedAt:      fixedTime,
						})
				},
			},
			expected: Expected{
				education: nil,
				err:       fmt.Errorf("invalid education returned: level invalid = invalid-level"),
			},
		},
		"Missing createdAt": {
			given: Given{
				id: fixedID,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(mock.Anything, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
						return len(args) == 1 && args[0] == fixedID
					})).
						Return(&educationGetFakeRow{
							id:             fixedID,
							mainSchoolJSON: mainSchoolJSON,
							level:          domain.College,
							updatedAt:      fixedTime,
						})
				},
			},
			expected: Expected{
				education: nil,
				err:       fmt.Errorf("invalid education returned: createdAt missing"),
			},
		},
		"Missing updatedAt": {
			given: Given{
				id: fixedID,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(mock.Anything, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
						return len(args) == 1 && args[0] == fixedID
					})).
						Return(&educationGetFakeRow{
							id:             fixedID,
							mainSchoolJSON: mainSchoolJSON,
							level:          domain.College,
							createdAt:      fixedTime,
						})
				},
			},
			expected: Expected{
				education: nil,
				err:       fmt.Errorf("invalid education returned: updatedAt missing"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newEducationRepositoryTestFixture(t, func() time.Time { return fixedTime })

			if test.given.mockQueryRow != nil {
				test.given.mockQueryRow(f.databaseAPI)
			}

			edu, err := f.educationRepository.Get(context.Background(), test.given.id)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
				assert.Nil(t, edu)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected.education, edu)
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}

func TestEducationRepository_Update(t *testing.T) {
	fixedID := "123-abc"
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	scanErr := errors.New("scan error")

	validMainSchool := domain.SchoolPeriod{
		Link:        "http://example.com",
		Name:        "test-name",
		Description: "test-description",
		Logo:        "test-logo",
		BlurHash:    "test-blurhash",
		Honor:       "test-honor",
		StartDate:   time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	validProjects := []domain.Project{
		{
			Preview:     "test-preview",
			BlurHash:    "test-blurhash",
			Title:       "test-title",
			SubTitle:    "test-subtitle",
			Description: "test-description",
			Stack:       []string{"stack1"},
			Type:        domain.Web,
			Link:        "http://example.com",
		},
	}

	validEducation := domain.Education{
		Id:            fixedID,
		MainSchool:    validMainSchool,
		SchoolPeriods: []domain.SchoolPeriod{validMainSchool},
		Projects:      validProjects,
		Level:         domain.College,
		CreatedAt:     fixedTime,
		UpdatedAt:     fixedTime,
	}

	mainSchoolJSON, _ := json.Marshal(validMainSchool)
	schoolPeriodsJSON, _ := json.Marshal([]domain.SchoolPeriod{validMainSchool})
	projectsJSON, _ := json.Marshal(validProjects)

	type Given struct {
		education    domain.Education
		mockQueryRow func(m *database.MockDatabaseAPI)
	}

	type Expected struct {
		education *domain.Education
		err       error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful update education": {
			given: Given{
				education: validEducation,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(
							mock.Anything,
							mock.Anything,
							mock.MatchedBy(func(args []interface{}) bool {
								return len(args) == 6 &&
									args[0] == fixedID &&
									bytes.Equal(args[1].([]byte), mainSchoolJSON)
							}),
						).
						Return(&educationGetFakeRow{
							id:                fixedID,
							mainSchoolJSON:    mainSchoolJSON,
							schoolPeriodsJSON: schoolPeriodsJSON,
							projectsJSON:      projectsJSON,
							level:             domain.College,
							createdAt:         fixedTime,
							updatedAt:         fixedTime,
						})
				},
			},
			expected: Expected{
				education: &domain.Education{
					Id:            fixedID,
					MainSchool:    validMainSchool,
					SchoolPeriods: []domain.SchoolPeriod{validMainSchool},
					Projects:      validProjects,
					Level:         domain.College,
					CreatedAt:     fixedTime,
					UpdatedAt:     fixedTime,
				},
				err: nil,
			},
		},
		"Database returns no rows": {
			given: Given{
				education: validEducation,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, mock.Anything).
						Return(&educationGetFakeRow{scanErr: pgx.ErrNoRows})
				},
			},
			expected: Expected{
				education: nil,
				err:       nil,
			},
		},
		"Database scan fails": {
			given: Given{
				education: validEducation,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, mock.Anything).
						Return(&educationGetFakeRow{scanErr: scanErr})
				},
			},
			expected: Expected{
				education: nil,
				err:       fmt.Errorf("failed to update education: %w", scanErr),
			},
		},
		"Invalid JSON returned": {
			given: Given{
				education: validEducation,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, mock.Anything).
						Return(&educationGetFakeRow{
							id:             fixedID,
							mainSchoolJSON: []byte("{invalid-json}"),
							level:          domain.College,
							createdAt:      fixedTime,
							updatedAt:      fixedTime,
						})
				},
			},
			expected: Expected{
				education: nil,
				err:       fmt.Errorf("failed to unmarshal main school: %w", errors.New("invalid character 'i' looking for beginning of object key string")),
			},
		},
		"Missing ID fails": {
			given: Given{
				education: domain.Education{
					MainSchool:    validEducation.MainSchool,
					SchoolPeriods: validEducation.SchoolPeriods,
					Projects:      validEducation.Projects,
					Level:         validEducation.Level,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to update education: ID missing"),
			},
		},
		"Missing main school fails": {
			given: Given{
				education: domain.Education{
					Id:            validEducation.Id,
					MainSchool:    domain.SchoolPeriod{},
					SchoolPeriods: validEducation.SchoolPeriods,
					Projects:      validEducation.Projects,
					Level:         validEducation.Level,
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
					Id: validEducation.Id,
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
					Id: validEducation.Id,
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
					Id: validEducation.Id,
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
					Id: validEducation.Id,
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
					Id: validEducation.Id,
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
					Id: validEducation.Id,
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
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: main school end date missing"),
			},
		},
		"End date before start date fails": {
			given: Given{
				education: domain.Education{
					Id: validEducation.Id,
					MainSchool: domain.SchoolPeriod{
						Name:        validEducation.MainSchool.Name,
						Description: validEducation.MainSchool.Description,
						Logo:        validEducation.MainSchool.Logo,
						BlurHash:    validEducation.MainSchool.BlurHash,
						StartDate:   time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
						EndDate:     time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC),
					},
					SchoolPeriods: validEducation.SchoolPeriods,
					Projects:      validEducation.Projects,
					Level:         validEducation.Level,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: main school end date must be after start date"),
			},
		},
		"Missing school period fails": {
			given: Given{
				education: domain.Education{
					Id:         validEducation.Id,
					MainSchool: validEducation.MainSchool,
					SchoolPeriods: []domain.SchoolPeriod{
						{},
					},
					Projects: validEducation.Projects,
					Level:    validEducation.Level,
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
					Id:         validEducation.Id,
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
					Projects: validEducation.Projects,
					Level:    validEducation.Level,
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
					Id:         validEducation.Id,
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
					Level: validEducation.Level,
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
					Id:         validEducation.Id,
					MainSchool: validEducation.MainSchool,
					Projects:   validEducation.Projects,
					Level:      "",
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
					Id:         validEducation.Id,
					MainSchool: validEducation.MainSchool,
					Projects:   validEducation.Projects,
					Level:      "invalid",
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate education: level invalid = invalid"),
			},
		},
		"Update with empty optional fields": {
			given: Given{
				education: domain.Education{
					Id:         fixedID,
					MainSchool: validMainSchool,
					Level:      domain.College,
					CreatedAt:  fixedTime,
					UpdatedAt:  fixedTime,
				},
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, mock.Anything).
						Return(&educationGetFakeRow{
							id:                fixedID,
							mainSchoolJSON:    mainSchoolJSON,
							schoolPeriodsJSON: []byte{},
							projectsJSON:      []byte{},
							level:             domain.College,
							createdAt:         fixedTime,
							updatedAt:         fixedTime,
						})
				},
			},
			expected: Expected{
				education: &domain.Education{
					Id:         fixedID,
					MainSchool: validMainSchool,
					Level:      domain.College,
					CreatedAt:  fixedTime,
					UpdatedAt:  fixedTime,
				},
				err: nil,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newEducationRepositoryTestFixture(t, func() time.Time { return fixedTime })

			if test.given.mockQueryRow != nil {
				test.given.mockQueryRow(f.databaseAPI)
			}

			edu, err := f.educationRepository.Update(context.Background(), &test.given.education)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
				assert.Nil(t, edu)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected.education, edu)
				assert.Equal(t, fixedTime, test.given.education.UpdatedAt)
				if edu != nil && test.expected.education != nil {
					assert.Equal(t, fixedTime, edu.CreatedAt, "CreatedAt should not change on update")
				}
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}

func TestEducationRepository_Delete(t *testing.T) {
	fixedId := "123-abc"
	dbErr := errors.New("db exec error")

	type Given struct {
		id       string
		mockExec func(m *database.MockDatabaseAPI)
	}

	type Expected struct {
		err error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful delete": {
			given: Given{
				id: fixedId,
				mockExec: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Exec(mock.Anything, mock.Anything, mock.Anything).
						Return(educationFakeCommandTag("DELETE 1"), nil)
				},
			},
			expected: Expected{err: nil},
		},
		"Database error during delete": {
			given: Given{
				id: fixedId,
				mockExec: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Exec(mock.Anything, mock.Anything, mock.Anything).
						Return(nil, dbErr)
				},
			},
			expected: Expected{
				err: fmt.Errorf("failed to delete education: %w", dbErr),
			},
		},
		"No rows affected": {
			given: Given{
				id: fixedId,
				mockExec: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Exec(mock.Anything, mock.Anything, mock.Anything).
						Return(educationFakeCommandTag("DELETE 0"), nil)
				},
			},
			expected: Expected{
				err: pgx.ErrNoRows,
			},
		},
		"Delete with empty ID": {
			given: Given{
				id:       "",
				mockExec: nil,
			},
			expected: Expected{
				err: fmt.Errorf("failed to delete education: ID missing"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newEducationRepositoryTestFixture(t, time.Now)

			if test.given.mockExec != nil {
				test.given.mockExec(f.databaseAPI)
			}

			err := f.educationRepository.Delete(context.Background(), test.given.id)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
			} else {
				assert.NoError(t, err)
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}

func TestEducationRepository_List(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	scanErr := errors.New("scan error")
	queryErr := errors.New("query error")
	rowErr := errors.New("row iteration error")

	validMainSchool := domain.SchoolPeriod{
		Link:        "http://example.com",
		Name:        "test-name",
		Description: "test-description",
		Logo:        "test-logo",
		BlurHash:    "test-blurhash",
		Honor:       "test-honor",
		StartDate:   time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	validProjects := []domain.Project{
		{
			Preview:     "test-preview",
			BlurHash:    "test-blurhash",
			Title:       "test-title",
			SubTitle:    "test-subtitle",
			Description: "test-description",
			Stack:       []string{"stack1"},
			Type:        domain.Web,
			Link:        "http://example.com",
		},
	}

	validEducation := domain.Education{
		Id:         "edu-001",
		MainSchool: validMainSchool,
		Projects:   validProjects,
		Level:      "Bachelor",
		CreatedAt:  fixedTime,
		UpdatedAt:  fixedTime,
	}

	type Given struct {
		filter    domain.EducationFilter
		mockQuery func(m *database.MockDatabaseAPI)
	}

	type Expected struct {
		education []domain.Education
		err       error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful list education": {
			given: Given{
				filter: domain.EducationFilter{},
				mockQuery: func(m *database.MockDatabaseAPI) {
					rows := &educationFakeRows{
						rows: []*educationFakeRow{
							{education: validEducation},
						},
					}
					m.EXPECT().
						Query(mock.Anything, mock.Anything).
						Return(rows, nil)
				},
			},
			expected: Expected{
				education: []domain.Education{validEducation},
				err:       nil,
			},
		},
		"Query fails": {
			given: Given{
				filter: domain.EducationFilter{},
				mockQuery: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Query(mock.Anything, mock.Anything).
						Return(nil, queryErr)
				},
			},
			expected: Expected{
				education: nil,
				err:       fmt.Errorf("failed to list education: %w", queryErr),
			},
		},
		"Scan fails": {
			given: Given{
				filter: domain.EducationFilter{},
				mockQuery: func(m *database.MockDatabaseAPI) {
					rows := &educationFakeRows{
						rows: []*educationFakeRow{
							{scanErr: scanErr},
						},
					}
					m.EXPECT().
						Query(mock.Anything, mock.Anything).
						Return(rows, nil)
				},
			},
			expected: Expected{
				education: nil,
				err:       fmt.Errorf("failed to scan education: %w", scanErr),
			},
		},
		"Row iteration error": {
			given: Given{
				filter: domain.EducationFilter{},
				mockQuery: func(m *database.MockDatabaseAPI) {
					rows := &educationFakeRows{
						rows:   []*educationFakeRow{},
						rowErr: rowErr,
					}
					m.EXPECT().
						Query(mock.Anything, mock.Anything).
						Return(rows, nil)
				},
			},
			expected: Expected{
				education: nil,
				err:       fmt.Errorf("row iteration error: %w", rowErr),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newEducationRepositoryTestFixture(t, func() time.Time { return fixedTime })

			if test.given.mockQuery != nil {
				test.given.mockQuery(f.databaseAPI)
			}

			education, err := f.educationRepository.List(context.Background(), test.given.filter)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
				assert.Nil(t, education)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, education)
				assert.Equal(t, test.expected.education, education)
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}
