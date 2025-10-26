package v1

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	database "github.com/fingertips18/fingertips18.github.io/backend/internal/database/mocks"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testSkillTable = "test-skills"
)

// skillFakeRow is for QueryRow
type skillFakeRow struct {
	id      string
	scanErr error
}

func (f *skillFakeRow) Scan(dest ...any) error {
	if f.scanErr != nil {
		return f.scanErr
	}

	if len(dest) != 1 {
		return fmt.Errorf("expected 1 scan destination, got %d", len(dest))
	}

	if v, ok := dest[0].(*string); ok {
		*v = f.id
		return nil
	}

	return fmt.Errorf("expected *string for id, got %T", dest[0])
}

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

	validSkill := domain.CreateSkill{
		Icon:     "🔥",
		HexColor: "#FF0000",
		Label:    "Go",
		Category: "Backend",
	}

	type Given struct {
		skill        domain.CreateSkill
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
				skill: domain.CreateSkill{
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
				skill: domain.CreateSkill{
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
				skill: domain.CreateSkill{
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
				skill: domain.CreateSkill{
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
