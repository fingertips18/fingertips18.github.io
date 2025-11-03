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
	testProjectTable = "test-projects"
)

// projectFakeRow is for QueryRow
type projectFakeRow struct {
	id      string
	project domain.Project
	scanErr error
}

func (f *projectFakeRow) Scan(dest ...any) error {
	if f.scanErr != nil {
		return f.scanErr
	}

	switch len(dest) {
	case 1:
		if v, ok := dest[0].(*string); ok {
			*v = f.id
			return nil
		}
		return fmt.Errorf("expected *string for id, got %T", dest[0])

	case 11:
		*dest[0].(*string) = f.project.Id
		*dest[1].(*string) = f.project.Preview
		*dest[2].(*string) = f.project.BlurHash
		*dest[3].(*string) = f.project.Title
		*dest[4].(*string) = f.project.SubTitle
		*dest[5].(*string) = f.project.Description

		if v, ok := dest[6].(*[]string); ok {
			*v = f.project.Stack
		} else {
			return fmt.Errorf("expected *[]string for stack, got %T", dest[6])
		}

		switch v := dest[7].(type) {
		case *string:
			*v = string(f.project.Type)
		case *domain.ProjectType:
			*v = f.project.Type
		default:
			return fmt.Errorf("unexpected type for project.Type: %T", v)
		}

		*dest[8].(*string) = f.project.Link
		*dest[9].(*time.Time) = f.project.CreatedAt
		*dest[10].(*time.Time) = f.project.UpdatedAt
		return nil

	default:
		return fmt.Errorf("unsupported number of scan destinations: %d", len(dest))
	}
}

type projectFakeCommandTag string

func (f projectFakeCommandTag) RowsAffected() int64 {
	if f == "DELETE 1" {
		return 1
	}
	return 0
}

// projectFakeRows is for Query
type projectFakeRows struct {
	rows   []*projectFakeRow
	index  int
	rowErr error
}

func (r *projectFakeRows) Next() bool {
	return r.index < len(r.rows)
}

func (r *projectFakeRows) Scan(dest ...any) error {
	if r.index >= len(r.rows) {
		return io.EOF
	}
	err := r.rows[r.index].Scan(dest...)
	r.index++
	return err
}

func (r *projectFakeRows) Err() error { return r.rowErr }

// Must match database.Rows interface (no return)
func (r *projectFakeRows) Close() {}

type projectRepositoryTestFixture struct {
	t                 *testing.T
	databaseAPI       *database.MockDatabaseAPI
	projectRepository *projectRepository
}

func newProjectRepositoryTestFixture(t *testing.T, timeProvider func() time.Time) *projectRepositoryTestFixture {
	mockDatabaseAPI := new(database.MockDatabaseAPI)
	projectRepository := &projectRepository{
		databaseAPI:  mockDatabaseAPI,
		timeProvider: timeProvider,
		projectTable: testProjectTable,
	}

	return &projectRepositoryTestFixture{
		t:                 t,
		databaseAPI:       mockDatabaseAPI,
		projectRepository: projectRepository,
	}
}

func TestProjectRepository_Create(t *testing.T) {
	fixedID := "123-abc"
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	scanErr := errors.New("scan error")

	validProject := domain.Project{
		Preview:     "test-preview",
		BlurHash:    "test-blurhash",
		Title:       "test-title",
		SubTitle:    "test-subtitle",
		Description: "test-description",
		Stack:       []string{"stack1"},
		Type:        domain.Web,
		Link:        "http://example.com",
	}

	type Given struct {
		project      domain.Project
		mockQueryRow func(m *database.MockDatabaseAPI)
	}

	type Expected struct {
		err error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful create project": {
			given: Given{
				project: validProject,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.MatchedBy(func(query string) bool { return strings.Contains(query, "INSERT INTO") }),
						mock.AnythingOfType("[]interface {}"),
					).Return(&projectFakeRow{id: fixedID})
				},
			},
			expected: Expected{
				err: nil,
			},
		},
		"Database scan fails": {
			given: Given{
				project: validProject,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.MatchedBy(func(query string) bool { return strings.Contains(query, "INSERT INTO") }),
						mock.AnythingOfType("[]interface {}"),
					).Return(&projectFakeRow{scanErr: scanErr})
				},
			},
			expected: Expected{
				err: fmt.Errorf("failed to create project: %w", scanErr),
			},
		},
		"Database returns empty ID": {
			given: Given{
				project: validProject,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(
						mock.Anything,
						mock.MatchedBy(func(query string) bool { return strings.Contains(query, "INSERT INTO") }),
						mock.AnythingOfType("[]interface {}"),
					).Return(&projectFakeRow{id: ""})
				},
			},
			expected: Expected{
				err: errors.New("invalid project returned: ID missing"),
			},
		},
		"Missing preview fails": {
			given: Given{
				project: domain.Project{
					Preview:     "",
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       validProject.Stack,
					Type:        validProject.Type,
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate project: preview missing"),
			},
		},
		"Missing blurHash fails": {
			given: Given{
				project: domain.Project{
					Preview:     validProject.Preview,
					BlurHash:    "",
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       validProject.Stack,
					Type:        validProject.Type,
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate project: blurHash missing"),
			},
		},
		"Missing title fails": {
			given: Given{
				project: domain.Project{
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       "",
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       validProject.Stack,
					Type:        validProject.Type,
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate project: title missing"),
			},
		},
		"Missing subTitle fails": {
			given: Given{
				project: domain.Project{
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    "",
					Description: validProject.Description,
					Stack:       validProject.Stack,
					Type:        validProject.Type,
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate project: subTitle missing"),
			},
		},
		"Missing description fails": {
			given: Given{
				project: domain.Project{
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: "",
					Stack:       validProject.Stack,
					Type:        validProject.Type,
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate project: description missing"),
			},
		},
		"Missing stack fails": {
			given: Given{
				project: domain.Project{
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       []string{},
					Type:        validProject.Type,
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate project: stack missing"),
			},
		},
		"Stack contains empty string fails": {
			given: Given{
				project: domain.Project{
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       []string{"stack1", "", "stack3"},
					Type:        validProject.Type,
					Link:        validProject.Link,
					CreatedAt:   fixedTime,
					UpdatedAt:   fixedTime,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate project: stack[1] is empty"),
			},
		},
		"Missing type fails": {
			given: Given{
				project: domain.Project{
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       validProject.Stack,
					Type:        "",
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate project: type missing"),
			},
		},
		"Invalid type fails": {
			given: Given{
				project: domain.Project{
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       validProject.Stack,
					Type:        "invalid",
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate project: type invalid = invalid"),
			},
		},
		"Missing link fails": {
			given: Given{
				project: domain.Project{
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       validProject.Stack,
					Type:        validProject.Type,
					Link:        "",
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				err: errors.New("failed to validate project: link missing"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newProjectRepositoryTestFixture(t, func() time.Time { return fixedTime })

			if test.given.mockQueryRow != nil {
				test.given.mockQueryRow(f.databaseAPI)
			}

			id, err := f.projectRepository.Create(context.Background(), &test.given.project)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, fixedID, id)
				assert.Equal(t, fixedTime, test.given.project.CreatedAt)
				assert.Equal(t, fixedTime, test.given.project.UpdatedAt)
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}

func TestProjectRepository_Get(t *testing.T) {
	id := "123-abc"
	scanErr := errors.New("scan error")

	validProject := domain.Project{
		Id:          id,
		Preview:     "test-preview",
		BlurHash:    "test-blurhash",
		Title:       "test-title",
		SubTitle:    "test-subtitle",
		Description: "test-description",
		Stack:       []string{"stack1"},
		Type:        domain.Web,
		Link:        "http://example.com",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	type Given struct {
		id           string
		mockQueryRow func(m *database.MockDatabaseAPI)
	}

	type Expected struct {
		project *domain.Project
		err     error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful get project": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{project: validProject})
				},
			},
			expected: Expected{
				project: &validProject,
				err:     nil,
			},
		},
		"Database scan": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{scanErr: scanErr})
				},
			},
			expected: Expected{
				project: nil,
				err:     fmt.Errorf("failed to scan project: %w", scanErr),
			},
		},
		"Database no rows": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{scanErr: pgx.ErrNoRows})
				},
			},
			expected: Expected{
				project: nil,
				err:     fmt.Errorf("failed to get project: %w", pgx.ErrNoRows),
			},
		},
		"Missing ID fails": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{project: domain.Project{
							Id:          "",
							Preview:     validProject.Preview,
							BlurHash:    validProject.BlurHash,
							Title:       validProject.Title,
							SubTitle:    validProject.SubTitle,
							Description: validProject.Description,
							Stack:       validProject.Stack,
							Type:        validProject.Type,
							Link:        validProject.Link,
							CreatedAt:   validProject.CreatedAt,
							UpdatedAt:   validProject.UpdatedAt,
						}})
				},
			},
			expected: Expected{
				project: nil,
				err:     errors.New("invalid project returned: ID missing"),
			},
		},
		"Missing preview fails": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{project: domain.Project{
							Id:          validProject.Id,
							Preview:     "",
							BlurHash:    validProject.BlurHash,
							Title:       validProject.Title,
							SubTitle:    validProject.SubTitle,
							Description: validProject.Description,
							Stack:       validProject.Stack,
							Type:        validProject.Type,
							Link:        validProject.Link,
							CreatedAt:   validProject.CreatedAt,
							UpdatedAt:   validProject.UpdatedAt,
						}})
				},
			},
			expected: Expected{
				project: nil,
				err:     errors.New("invalid project returned: preview missing"),
			},
		},
		"Missing blurHash fails": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{project: domain.Project{
							Id:          validProject.Id,
							Preview:     validProject.Preview,
							BlurHash:    "",
							Title:       validProject.Title,
							SubTitle:    validProject.SubTitle,
							Description: validProject.Description,
							Stack:       validProject.Stack,
							Type:        validProject.Type,
							Link:        validProject.Link,
							CreatedAt:   validProject.CreatedAt,
							UpdatedAt:   validProject.UpdatedAt,
						}})
				},
			},
			expected: Expected{
				project: nil,
				err:     errors.New("invalid project returned: blurHash missing"),
			},
		},
		"Missing title fails": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{project: domain.Project{
							Id:          validProject.Id,
							Preview:     validProject.Preview,
							BlurHash:    validProject.BlurHash,
							Title:       "",
							SubTitle:    validProject.SubTitle,
							Description: validProject.Description,
							Stack:       validProject.Stack,
							Type:        validProject.Type,
							Link:        validProject.Link,
							CreatedAt:   validProject.CreatedAt,
							UpdatedAt:   validProject.UpdatedAt,
						}})
				},
			},
			expected: Expected{
				project: nil,
				err:     errors.New("invalid project returned: title missing"),
			},
		},
		"Missing subTitle fails": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{project: domain.Project{
							Id:          validProject.Id,
							Preview:     validProject.Preview,
							BlurHash:    validProject.BlurHash,
							Title:       validProject.Title,
							SubTitle:    "",
							Description: validProject.Description,
							Stack:       validProject.Stack,
							Type:        validProject.Type,
							Link:        validProject.Link,
							CreatedAt:   validProject.CreatedAt,
							UpdatedAt:   validProject.UpdatedAt,
						}})
				},
			},
			expected: Expected{
				project: nil,
				err:     errors.New("invalid project returned: subTitle missing"),
			},
		},
		"Missing description fails": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{project: domain.Project{
							Id:          validProject.Id,
							Preview:     validProject.Preview,
							BlurHash:    validProject.BlurHash,
							Title:       validProject.Title,
							SubTitle:    validProject.SubTitle,
							Description: "",
							Stack:       validProject.Stack,
							Type:        validProject.Type,
							Link:        validProject.Link,
							CreatedAt:   validProject.CreatedAt,
							UpdatedAt:   validProject.UpdatedAt,
						}})
				},
			},
			expected: Expected{
				project: nil,
				err:     errors.New("invalid project returned: description missing"),
			},
		},
		"Missing stack fails": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{project: domain.Project{
							Id:          validProject.Id,
							Preview:     validProject.Preview,
							BlurHash:    validProject.BlurHash,
							Title:       validProject.Title,
							SubTitle:    validProject.SubTitle,
							Description: validProject.Description,
							Stack:       nil,
							Type:        validProject.Type,
							Link:        validProject.Link,
							CreatedAt:   validProject.CreatedAt,
							UpdatedAt:   validProject.UpdatedAt,
						}})
				},
			},
			expected: Expected{
				project: nil,
				err:     errors.New("invalid project returned: stack missing"),
			},
		},
		"Missing type fails": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{project: domain.Project{
							Id:          validProject.Id,
							Preview:     validProject.Preview,
							BlurHash:    validProject.BlurHash,
							Title:       validProject.Title,
							SubTitle:    validProject.SubTitle,
							Description: validProject.Description,
							Stack:       validProject.Stack,
							Type:        "",
							Link:        validProject.Link,
							CreatedAt:   validProject.CreatedAt,
							UpdatedAt:   validProject.UpdatedAt,
						}})
				},
			},
			expected: Expected{
				project: nil,
				err:     errors.New("invalid project returned: type missing"),
			},
		},
		"Missing link fails": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{project: domain.Project{
							Id:          validProject.Id,
							Preview:     validProject.Preview,
							BlurHash:    validProject.BlurHash,
							Title:       validProject.Title,
							SubTitle:    validProject.SubTitle,
							Description: validProject.Description,
							Stack:       validProject.Stack,
							Type:        validProject.Type,
							Link:        "",
							CreatedAt:   validProject.CreatedAt,
							UpdatedAt:   validProject.UpdatedAt,
						}})
				},
			},
			expected: Expected{
				project: nil,
				err:     errors.New("invalid project returned: link missing"),
			},
		},
		"Missing createdAt fails": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{project: domain.Project{
							Id:          validProject.Id,
							Preview:     validProject.Preview,
							BlurHash:    validProject.BlurHash,
							Title:       validProject.Title,
							SubTitle:    validProject.SubTitle,
							Description: validProject.Description,
							Stack:       validProject.Stack,
							Type:        validProject.Type,
							Link:        validProject.Link,
							CreatedAt:   time.Time{},
							UpdatedAt:   validProject.UpdatedAt,
						}})
				},
			},
			expected: Expected{
				project: nil,
				err:     errors.New("invalid project returned: createdAt missing"),
			},
		},
		"Missing updatedAt fails": {
			given: Given{
				id: id,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, []any{id}).
						Return(&projectFakeRow{project: domain.Project{
							Id:          validProject.Id,
							Preview:     validProject.Preview,
							BlurHash:    validProject.BlurHash,
							Title:       validProject.Title,
							SubTitle:    validProject.SubTitle,
							Description: validProject.Description,
							Stack:       validProject.Stack,
							Type:        validProject.Type,
							Link:        validProject.Link,
							CreatedAt:   validProject.CreatedAt,
							UpdatedAt:   time.Time{},
						}})
				},
			},
			expected: Expected{
				project: nil,
				err:     errors.New("invalid project returned: updatedAt missing"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newProjectRepositoryTestFixture(t, time.Now)

			if test.given.mockQueryRow != nil {
				test.given.mockQueryRow(f.databaseAPI)
			}

			project, err := f.projectRepository.Get(context.Background(), test.given.id)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Empty(t, err)
				assert.NotEmpty(t, project)
				assert.Equal(t, test.expected.project.Id, project.Id)
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}

func TestProjectRepository_Update(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	scanErr := errors.New("scan error")

	validProject := &domain.Project{
		Id:          "123-abc",
		Preview:     "test-preview",
		BlurHash:    "test-blurhash",
		Title:       "test-title",
		SubTitle:    "test-subtitle",
		Description: "test-description",
		Stack:       []string{"stack1"},
		Type:        domain.Web,
		Link:        "http://example.com",
	}

	type Given struct {
		project      domain.Project
		mockQueryRow func(m *database.MockDatabaseAPI)
	}

	type Expected struct {
		updatedProject *domain.Project
		err            error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful update project": {
			given: Given{
				project: *validProject,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, mock.Anything).
						Return(&projectFakeRow{
							project: *validProject,
						})
				},
			},
			expected: Expected{
				updatedProject: validProject,
				err:            nil,
			},
		},
		"Database scan fails": {
			given: Given{
				project: *validProject,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, mock.Anything).
						Return(&projectFakeRow{scanErr: scanErr})
				},
			},
			expected: Expected{
				updatedProject: nil,
				err:            fmt.Errorf("failed to update project: %w", scanErr),
			},
		},
		"Database returns no rows": {
			given: Given{
				project: *validProject,
				mockQueryRow: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						QueryRow(mock.Anything, mock.Anything, mock.Anything).
						Return(&projectFakeRow{scanErr: pgx.ErrNoRows})
				},
			},
			expected: Expected{
				updatedProject: nil,
				err:            nil,
			},
		},
		"Missing preview fails": {
			given: Given{
				project: domain.Project{
					Id:          validProject.Id,
					Preview:     "",
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       validProject.Stack,
					Type:        validProject.Type,
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				updatedProject: nil,
				err:            errors.New("failed to validate project: preview missing"),
			},
		},
		"Missing blurHash fails": {
			given: Given{
				project: domain.Project{
					Id:          validProject.Id,
					Preview:     validProject.Preview,
					BlurHash:    "",
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       validProject.Stack,
					Type:        validProject.Type,
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				updatedProject: nil,
				err:            errors.New("failed to validate project: blurHash missing"),
			},
		},
		"Missing title fails": {
			given: Given{
				project: domain.Project{
					Id:          validProject.Id,
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       "",
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       validProject.Stack,
					Type:        validProject.Type,
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				updatedProject: nil,
				err:            errors.New("failed to validate project: title missing"),
			},
		},
		"Missing subTitle fails": {
			given: Given{
				project: domain.Project{
					Id:          validProject.Id,
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    "",
					Description: validProject.Description,
					Stack:       validProject.Stack,
					Type:        validProject.Type,
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				updatedProject: nil,
				err:            errors.New("failed to validate project: subTitle missing"),
			},
		},
		"Missing description fails": {
			given: Given{
				project: domain.Project{
					Id:          validProject.Id,
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: "",
					Stack:       validProject.Stack,
					Type:        validProject.Type,
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				updatedProject: nil,
				err:            errors.New("failed to validate project: description missing"),
			},
		},
		"Missing stack fails": {
			given: Given{
				project: domain.Project{
					Id:          validProject.Id,
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       []string{},
					Type:        validProject.Type,
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				updatedProject: nil,
				err:            errors.New("failed to validate project: stack missing"),
			},
		},
		"Stack contains empty string fails": {
			given: Given{
				project: domain.Project{
					Id:          validProject.Id,
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       []string{"stack1", "", "stack3"},
					Type:        validProject.Type,
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				updatedProject: nil,
				err:            errors.New("failed to validate project: stack[1] is empty"),
			},
		},
		"Missing type fails": {
			given: Given{
				project: domain.Project{
					Id:          validProject.Id,
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       validProject.Stack,
					Type:        "",
					Link:        validProject.Link,
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				updatedProject: nil,
				err:            errors.New("failed to validate project: type missing"),
			},
		},
		"Missing link fails": {
			given: Given{
				project: domain.Project{
					Id:          validProject.Id,
					Preview:     validProject.Preview,
					BlurHash:    validProject.BlurHash,
					Title:       validProject.Title,
					SubTitle:    validProject.SubTitle,
					Description: validProject.Description,
					Stack:       validProject.Stack,
					Type:        validProject.Type,
					Link:        "",
				},
				mockQueryRow: nil,
			},
			expected: Expected{
				updatedProject: nil,
				err:            errors.New("failed to validate project: link missing"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newProjectRepositoryTestFixture(t, func() time.Time { return fixedTime })

			if test.given.mockQueryRow != nil {
				test.given.mockQueryRow(f.databaseAPI)
			}

			project, err := f.projectRepository.Update(context.Background(), &test.given.project)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
				assert.Nil(t, project)
			} else {
				assert.NoError(t, err)
				if test.expected.updatedProject == nil {
					assert.Nil(t, project)
				} else {
					assert.NotNil(t, project)
					assert.Equal(t, test.expected.updatedProject.Id, project.Id)
				}
				assert.Equal(t, fixedTime, test.given.project.UpdatedAt)
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}

func TestProjectRepository_Delete(t *testing.T) {
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
				id: "123-abc",
				mockExec: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Exec(mock.Anything, mock.Anything, mock.Anything).
						Return(projectFakeCommandTag("DELETE 1"), nil)
				},
			},
			expected: Expected{err: nil},
		},
		"Database error during delete": {
			given: Given{
				id: "123-abc",
				mockExec: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Exec(mock.Anything, mock.Anything, mock.Anything).
						Return(nil, dbErr)
				},
			},
			expected: Expected{
				err: fmt.Errorf("failed to delete project: %w", dbErr),
			},
		},
		"No rows affected": {
			given: Given{
				id: "123-abc",
				mockExec: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Exec(mock.Anything, mock.Anything, mock.Anything).
						Return(projectFakeCommandTag("DELETE 0"), nil)
				},
			},
			expected: Expected{
				err: pgx.ErrNoRows,
			},
		},
		"Delete with empty ID": {
			given: Given{
				id:       "",
				mockExec: nil, // no DB call expected
			},
			expected: Expected{
				err: fmt.Errorf("failed to delete project: ID missing"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newProjectRepositoryTestFixture(t, time.Now)

			if test.given.mockExec != nil {
				test.given.mockExec(f.databaseAPI)
			}

			err := f.projectRepository.Delete(context.Background(), test.given.id)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
			} else {
				assert.NoError(t, err)
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}

func TestProjectRepository_List(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	scanErr := errors.New("scan error")
	queryErr := errors.New("query error")
	rowErr := errors.New("row iteration error")

	validProject := domain.Project{
		Id:          "123-abc",
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
	}

	type Given struct {
		filter    domain.ProjectFilter
		mockQuery func(m *database.MockDatabaseAPI)
	}

	type Expected struct {
		projects []domain.Project
		err      error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful list projects": {
			given: Given{
				filter: domain.ProjectFilter{},
				mockQuery: func(m *database.MockDatabaseAPI) {
					rows := &projectFakeRows{
						rows: []*projectFakeRow{
							{project: validProject},
						},
					}
					m.EXPECT().
						Query(mock.Anything, mock.Anything, mock.Anything).
						Return(rows, nil)
				},
			},
			expected: Expected{
				projects: []domain.Project{validProject},
				err:      nil,
			},
		},
		"Query fails": {
			given: Given{
				filter: domain.ProjectFilter{},
				mockQuery: func(m *database.MockDatabaseAPI) {
					m.EXPECT().
						Query(mock.Anything, mock.Anything, mock.Anything).
						Return(nil, queryErr)
				},
			},
			expected: Expected{
				projects: nil,
				err:      fmt.Errorf("failed to list projects: %w", queryErr),
			},
		},
		"Scan fails": {
			given: Given{
				filter: domain.ProjectFilter{},
				mockQuery: func(m *database.MockDatabaseAPI) {
					rows := &projectFakeRows{
						rows: []*projectFakeRow{
							{scanErr: scanErr},
						},
					}
					m.EXPECT().
						Query(mock.Anything, mock.Anything, mock.Anything).
						Return(rows, nil)
				},
			},
			expected: Expected{
				projects: nil,
				err:      fmt.Errorf("failed to scan project: %w", scanErr),
			},
		},
		"Row iteration error": {
			given: Given{
				filter: domain.ProjectFilter{},
				mockQuery: func(m *database.MockDatabaseAPI) {
					rows := &projectFakeRows{
						rows:   []*projectFakeRow{},
						rowErr: rowErr,
					}
					m.EXPECT().
						Query(mock.Anything, mock.Anything, mock.Anything).
						Return(rows, nil)
				},
			},
			expected: Expected{
				projects: nil,
				err:      fmt.Errorf("row iteration error: %w", rowErr),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newProjectRepositoryTestFixture(t, func() time.Time { return fixedTime })

			if test.given.mockQuery != nil {
				test.given.mockQuery(f.databaseAPI)
			}

			projects, err := f.projectRepository.List(context.Background(), test.given.filter)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
				assert.Nil(t, projects)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, projects)
				assert.Equal(t, test.expected.projects, projects)
			}

			f.databaseAPI.AssertExpectations(t)
		})
	}
}
