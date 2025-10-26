package v1

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	database "github.com/fingertips18/fingertips18.github.io/backend/internal/database/mocks"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testSkillTable = "test-skills"
)

// skillFakeRow for Create, Get, and Update paths
type skillFakeRow struct {
	id      string
	skill   *domain.Skill
	scanErr error
}

func (f *skillFakeRow) Scan(dest ...any) error {
	if f.scanErr != nil {
		return f.scanErr
	}

	switch len(dest) {
	case 1: // Create() only returns ID
		*dest[0].(*string) = f.id
		return nil

	case 7: // Get() reads full object
		if f.skill == nil {
			return fmt.Errorf("skill not provided for scan")
		}
		*dest[0].(*string) = f.skill.Id
		*dest[1].(*string) = f.skill.Icon
		*dest[2].(*string) = f.skill.HexColor
		*dest[3].(*string) = f.skill.Label
		*dest[4].(*domain.SkillCategory) = f.skill.Category
		*dest[5].(*time.Time) = f.skill.CreatedAt
		*dest[6].(*time.Time) = f.skill.UpdatedAt
		return nil

	default:
		return fmt.Errorf("unexpected number of scan destinations: %d", len(dest))
	}
}

// skillFakeCommandTag satisfies the CommandTag interface for Exec mocking
type skillFakeCommandTag struct {
	rows int64
}

func (f *skillFakeCommandTag) RowsAffected() int64 {
	return f.rows
}

type skillFakeRows struct {
	rows   []*skillFakeRow
	index  int
	rowErr error
}

func (r *skillFakeRows) Next() bool {
	return r.index < len(r.rows)
}

func (r *skillFakeRows) Scan(dest ...any) error {
	if r.index >= len(r.rows) {
		return io.EOF
	}
	err := r.rows[r.index].Scan(dest...)
	r.index++
	return err
}

func (r *skillFakeRows) Err() error { return r.rowErr }

func (r *skillFakeRows) Close() {}

type skillRepositoryTestFixture struct {
	t               *testing.T
	databaseAPI     *database.MockDatabaseAPI
	skillRepository *skillRepository
}

func newSkillRepositoryTestFixture(t *testing.T, timeProvider func() time.Time) *skillRepositoryTestFixture {
	mockDatabaseAPI := new(database.MockDatabaseAPI)
	skillRepository := &skillRepository{
		databaseAPI:  mockDatabaseAPI,
		timeProvider: timeProvider,
		skillTable:   testSkillTable,
	}

	return &skillRepositoryTestFixture{
		t:               t,
		databaseAPI:     mockDatabaseAPI,
		skillRepository: skillRepository,
	}
}

func TestSkillRepository_Create(t *testing.T) {
	fixedID := "test-id"
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	scanErr := errors.New("scan error")

	validSkill := domain.Skill{
		Icon:     "🔥",
		HexColor: "#FF0000",
		Label:    "Go",
		Category: domain.Backend,
	}

	type Given struct {
		skill        domain.Skill
		mockQueryRow func(m *database.MockDatabaseAPI)
	}

	type Expected struct {
		err error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful create skill": {
			given: Given{
				skill: validSkill,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.MatchedBy(func(query string) bool { return strings.Contains(query, "INSERT INTO") }),
						mock.AnythingOfType("[]interface {}"),
					).Return(&skillFakeRow{id: fixedID})
				},
			},
			expected: Expected{
				err: nil,
			},
		},
		"Database scan fails": {
			given: Given{
				skill: validSkill,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).Return(&skillFakeRow{scanErr: scanErr})
				},
			},
			expected: Expected{
				err: fmt.Errorf("failed to create skill: %w", scanErr),
			},
		},
		"Database returns empty ID": {
			given: Given{
				skill: validSkill,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.Anything,
						mock.Anything,
					).Return(&skillFakeRow{id: ""})
				},
			},
			expected: Expected{
				err: errors.New("invalid skill returned: ID missing"),
			},
		},
		"Missing icon fails": {
			given: Given{
				skill: domain.Skill{
					Icon:     "",
					HexColor: validSkill.HexColor,
					Label:    validSkill.Label,
					Category: validSkill.Category,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate skill: icon missing"),
			},
		},
		"Missing hexColor fails": {
			given: Given{
				skill: domain.Skill{
					Icon:     validSkill.Icon,
					HexColor: "",
					Label:    validSkill.Label,
					Category: validSkill.Category,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate skill: hex color missing"),
			},
		},
		"Missing label fails": {
			given: Given{
				skill: domain.Skill{
					Icon:     validSkill.Icon,
					HexColor: validSkill.HexColor,
					Label:    "",
					Category: validSkill.Category,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate skill: label missing"),
			},
		},
		"Missing category fails": {
			given: Given{
				skill: domain.Skill{
					Icon:     validSkill.Icon,
					HexColor: validSkill.HexColor,
					Label:    validSkill.Label,
					Category: "",
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate skill: category missing"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newSkillRepositoryTestFixture(t, func() time.Time { return fixedTime })

			if test.given.mockQueryRow != nil {
				test.given.mockQueryRow(f.databaseAPI)
			}

			id, err := f.skillRepository.Create(context.Background(), &test.given.skill)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, fixedID, id)
				assert.Equal(t, fixedTime, test.given.skill.CreatedAt)
				assert.Equal(t, fixedTime, test.given.skill.UpdatedAt)
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}

func TestSkillRepository_Get(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	scanErr := errors.New("scan error")

	existingSkill := domain.Skill{
		Id:        "skill-123",
		Icon:      "🔥",
		HexColor:  "#FF0000",
		Label:     "Go",
		Category:  domain.Backend,
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}

	type Given struct {
		id           string
		mockQueryRow func(m *database.MockDatabaseAPI)
	}

	type Expected struct {
		result *domain.Skill
		err    error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful fetch skill": {
			given: Given{
				id: existingSkill.Id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.MatchedBy(func(q string) bool { return strings.Contains(q, "SELECT") }),
						mock.MatchedBy(func(args []any) bool {
							return len(args) == 1 && args[0] == existingSkill.Id
						}),
					).Return(&skillFakeRow{skill: &existingSkill})
				},
			},
			expected: Expected{
				result: &existingSkill,
				err:    nil,
			},
		},
		"Missing ID fails": {
			given: Given{
				id:           "",
				mockQueryRow: nil,
			},
			expected: Expected{
				result: nil,
				err:    fmt.Errorf("failed to get skill: ID missing"),
			},
		},
		"Skill not found returns error": {
			given: Given{
				id: existingSkill.Id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.MatchedBy(func(q string) bool { return strings.Contains(q, "SELECT") }),
						mock.MatchedBy(func(args []any) bool {
							return len(args) == 1 && args[0] == existingSkill.Id
						}),
					).Return(&skillFakeRow{scanErr: pgx.ErrNoRows})
				},
			},
			expected: Expected{
				result: nil,
				err:    fmt.Errorf("failed to get skill: %w", pgx.ErrNoRows),
			},
		},
		"Database scan fails": {
			given: Given{
				id: existingSkill.Id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.Anything,
						mock.AnythingOfType("[]interface {}"),
					).Return(&skillFakeRow{scanErr: scanErr})
				},
			},
			expected: Expected{
				result: nil,
				err:    fmt.Errorf("failed to scan skill: %w", scanErr),
			},
		},
		"Returned skill fails validation": {
			given: Given{
				id: existingSkill.Id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					invalid := existingSkill
					invalid.Label = "" // simulate validation failure
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.Anything,
						mock.MatchedBy(func(args []any) bool {
							return len(args) == 1 && args[0] == existingSkill.Id
						}),
					).Return(&skillFakeRow{skill: &invalid})
				},
			},
			expected: Expected{
				result: nil,
				err:    fmt.Errorf("invalid skill returned: %w", errors.New("label missing")),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newSkillRepositoryTestFixture(t, func() time.Time { return fixedTime })

			if test.given.mockQueryRow != nil {
				test.given.mockQueryRow(f.databaseAPI)
			}

			result, err := f.skillRepository.Get(context.Background(), test.given.id)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected.result, result)
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}

func TestSkillRepository_Update(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	scanErr := errors.New("scan error")

	originalSkill := domain.Skill{
		Id:        "existing-id",
		Icon:      "🔥",
		HexColor:  "#FF0000",
		Label:     "Go",
		Category:  domain.Backend,
		CreatedAt: fixedTime.Add(-time.Hour),
	}
	updatedSkill := originalSkill

	type Given struct {
		skill        *domain.Skill
		mockQueryRow func(m *database.MockDatabaseAPI)
	}

	type Expected struct {
		result *domain.Skill
		err    error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful update": {
			given: Given{
				skill: func() *domain.Skill { s := originalSkill; return &s }(),
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					updatedReturned := updatedSkill
					updatedReturned.UpdatedAt = fixedTime
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.MatchedBy(func(q string) bool { return strings.Contains(q, "UPDATE") }),
						mock.MatchedBy(func(args []any) bool {
							return len(args) == 6 &&
								args[0] == originalSkill.Id &&
								args[5] == fixedTime
						}),
					).Return(&skillFakeRow{skill: &updatedReturned})
				},
			},
			expected: Expected{
				result: &domain.Skill{
					Id:        originalSkill.Id,
					Icon:      originalSkill.Icon,
					HexColor:  originalSkill.HexColor,
					Label:     originalSkill.Label,
					Category:  originalSkill.Category,
					CreatedAt: originalSkill.CreatedAt,
					UpdatedAt: fixedTime,
				},
				err: nil,
			},
		},
		"Nil payload fails": {
			given: Given{
				skill:        nil,
				mockQueryRow: nil,
			},
			expected: Expected{
				result: nil,
				err:    errors.New("failed to validate skill: payload is nil"),
			},
		},
		"Missing ID fails": {
			given: Given{
				skill: &domain.Skill{
					Id:        "",
					Icon:      "🔥",
					HexColor:  "#FF0000",
					Label:     "Go",
					Category:  "Backend",
					CreatedAt: fixedTime,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				result: nil,
				err:    fmt.Errorf("failed to update skill: ID missing"),
			},
		},
		"Validation failure": {
			given: Given{
				skill: &domain.Skill{
					Id:        "existing-id",
					Icon:      "",
					HexColor:  "#FF0000",
					Label:     "Go",
					Category:  "Backend",
					CreatedAt: fixedTime,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				result: nil,
				err:    errors.New("failed to validate skill: icon missing"),
			},
		},
		"Skill not found returns nil result without error": {
			given: Given{
				skill: &originalSkill,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.Anything,
						mock.AnythingOfType("[]interface {}"),
					).Return(&skillFakeRow{scanErr: pgx.ErrNoRows})
				},
			},
			expected: Expected{
				result: nil,
				err:    nil,
			},
		},
		"Database scan fails": {
			given: Given{
				skill: &originalSkill,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.Anything,
						mock.AnythingOfType("[]interface {}"),
					).Return(&skillFakeRow{scanErr: scanErr})
				},
			},
			expected: Expected{
				result: nil,
				err:    fmt.Errorf("failed to update skill: %w", scanErr),
			},
		},
		"Returned skill fails validation": {
			given: Given{
				skill: func() *domain.Skill { s := originalSkill; return &s }(),
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					invalid := originalSkill
					invalid.Label = ""
					m.EXPECT().QueryRow(mock.Anything, mock.Anything, mock.Anything).
						Return(&skillFakeRow{skill: &invalid})
				},
			},
			expected: Expected{
				result: nil,
				err:    fmt.Errorf("invalid skill returned: %w", errors.New("label missing")),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// fresh fixture per test
			f := newSkillRepositoryTestFixture(t, func() time.Time { return fixedTime })

			if test.given.mockQueryRow != nil {
				test.given.mockQueryRow(f.databaseAPI)
			}

			result, err := f.skillRepository.Update(context.Background(), test.given.skill)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected.result, result)
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}

func TestSkillRepository_Delete(t *testing.T) {
	dbErr := errors.New("db exec error")
	fixedID := "skill-123"

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
				id: fixedID,
				mockExec: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Exec(
							mock.Anything,
							mock.MatchedBy(func(q string) bool {
								return strings.Contains(q, fmt.Sprintf("DELETE FROM %s", pgx.Identifier{testSkillTable}.Sanitize()))
							}),
							mock.Anything,
						).
						Return(&skillFakeCommandTag{rows: 1}, nil)
				},
			},
			expected: Expected{err: nil},
		},
		"Database error during delete": {
			given: Given{
				id: fixedID,
				mockExec: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Exec(mock.Anything, mock.Anything, mock.Anything).
						Return(nil, dbErr)
				},
			},
			expected: Expected{
				err: fmt.Errorf("failed to delete skill: %w", dbErr),
			},
		},
		"No rows affected": {
			given: Given{
				id: fixedID,
				mockExec: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Exec(mock.Anything, mock.Anything, mock.Anything).
						Return(&skillFakeCommandTag{rows: 0}, nil)
				},
			},
			expected: Expected{
				err: pgx.ErrNoRows,
			},
		},
		"Delete with empty ID": {
			given: Given{id: ""},
			expected: Expected{
				err: fmt.Errorf("failed to delete skill: ID missing"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newSkillRepositoryTestFixture(t, time.Now)

			if test.given.mockExec != nil {
				test.given.mockExec(f.databaseAPI)
			}

			err := f.skillRepository.Delete(context.Background(), test.given.id)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
			} else {
				assert.NoError(t, err)
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}

func TestSkillRepository_List(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	dbErr := errors.New("db error")
	scanErr := errors.New("scan error")
	iterErr := errors.New("row iteration error")

	category := domain.Backend

	mockSkill := domain.Skill{
		Id:        "skill-123",
		Icon:      "🔥",
		HexColor:  "#FF0000",
		Label:     "Go",
		Category:  domain.Backend,
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}

	type Given struct {
		filter   domain.SkillFilter
		mockRows func(m *database.MockDatabaseAPI)
	}

	type Expected struct {
		result []domain.Skill
		err    error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful list without filters": {
			given: Given{
				filter: domain.SkillFilter{},
				mockRows: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Query(mock.Anything, mock.Anything, mock.Anything).
						Return(&skillFakeRows{
							rows: []*skillFakeRow{{skill: &mockSkill}},
						}, nil)
				},
			},
			expected: Expected{
				result: []domain.Skill{mockSkill},
				err:    nil,
			},
		},
		"Successful list with category filter": {
			given: Given{
				filter: domain.SkillFilter{
					Category: &category,
				},
				mockRows: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Query(
							mock.Anything,
							mock.Anything,
							mock.MatchedBy(func(args ...any) bool {
								if len(args) != 1 {
									return false
								}
								switch v := args[0].(type) {
								case domain.SkillCategory:
									return v == category
								case string:
									return v == string(category)
								default:
									return false
								}
							}),
						).
						Return(&skillFakeRows{
							rows: []*skillFakeRow{{skill: &mockSkill}},
						}, nil)
				},
			},
			expected: Expected{
				result: []domain.Skill{mockSkill},
				err:    nil,
			},
		},
		"Database query error": {
			given: Given{
				filter: domain.SkillFilter{},
				mockRows: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Query(mock.Anything, mock.Anything, mock.Anything).
						Return(nil, dbErr)
				},
			},
			expected: Expected{
				result: nil,
				err:    fmt.Errorf("failed to list skills: %w", dbErr),
			},
		},
		"Scan failure during iteration": {
			given: Given{
				filter: domain.SkillFilter{},
				mockRows: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Query(mock.Anything, mock.Anything, mock.Anything).
						Return(&skillFakeRows{
							rows: []*skillFakeRow{{scanErr: scanErr}},
						}, nil)
				},
			},
			expected: Expected{
				result: nil,
				err:    fmt.Errorf("failed to scan skill: %w", scanErr),
			},
		},
		"Row iteration error": {
			given: Given{
				filter: domain.SkillFilter{},
				mockRows: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Query(mock.Anything, mock.Anything, mock.Anything).
						Return(&skillFakeRows{
							rows:   []*skillFakeRow{{skill: &mockSkill}},
							rowErr: iterErr,
						}, nil)
				},
			},
			expected: Expected{
				result: nil,
				err:    fmt.Errorf("row iteration error: %w", iterErr),
			},
		},
		"Pagination with boundaries": {
			given: Given{
				filter: domain.SkillFilter{
					Page:     -5,
					PageSize: 999,
				},
				mockRows: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Query(mock.Anything, mock.Anything, mock.Anything).
						Return(&skillFakeRows{
							rows: []*skillFakeRow{{skill: &mockSkill}},
						}, nil)
				},
			},
			expected: Expected{
				result: []domain.Skill{mockSkill},
				err:    nil,
			},
		},
		"Sorting applied correctly": {
			given: Given{
				filter: domain.SkillFilter{
					SortAscending: false,
				},
				mockRows: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Query(
							mock.Anything,
							mock.MatchedBy(func(q string) bool {
								return strings.Contains(q, "ORDER BY created_at DESC")
							}),
							mock.Anything,
						).
						Return(&skillFakeRows{
							rows: []*skillFakeRow{{skill: &mockSkill}},
						}, nil)
				},
			},
			expected: Expected{
				result: []domain.Skill{mockSkill},
				err:    nil,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newSkillRepositoryTestFixture(t, func() time.Time { return fixedTime })

			if test.given.mockRows != nil {
				test.given.mockRows(f.databaseAPI)
			}

			result, err := f.skillRepository.List(context.Background(), test.given.filter)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected.result, result)
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}
