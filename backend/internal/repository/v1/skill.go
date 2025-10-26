package v1

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/utils"
	"github.com/jackc/pgx/v5"
)

type SkillRepository interface {
	Create(ctx context.Context, skill *domain.Skill) (string, error)
	Get(ctx context.Context, id string) (*domain.Skill, error)
	Update(ctx context.Context, skill *domain.Skill) (*domain.Skill, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter domain.SkillFilter) ([]domain.Skill, error)
}

type SkillRepositoryConfig struct {
	DatabaseAPI database.DatabaseAPI
	SkillTable  string

	timeProvider domain.TimeProvider
}

type skillRepository struct {
	skillTable   string
	databaseAPI  database.DatabaseAPI
	timeProvider domain.TimeProvider
}

// NewSkillRepository creates and returns a SkillRepository configured with the
// values from the provided SkillRepositoryConfig. The repository will use
// cfg.SkillTable and cfg.DatabaseAPI, and it resolves cfg.timeProvider for
// time-related operations; if cfg.timeProvider is nil the repository defaults
// to time.Now.
func NewSkillRepository(cfg SkillRepositoryConfig) SkillRepository {
	timeProvider := cfg.timeProvider
	if timeProvider == nil {
		timeProvider = time.Now
	}

	return &skillRepository{
		skillTable:   cfg.SkillTable,
		databaseAPI:  cfg.DatabaseAPI,
		timeProvider: timeProvider,
	}
}

// Create inserts a new skill row into the configured skill table and returns its ID.
//
// The method performs the following steps:
//  1. Validates the provided Skill payload via skill.ValidatePayload().
//  2. Sets created_at and updated_at timestamps using the configured timeProvider.
//  3. Executes the INSERT and captures the RETURNING id clause.
//  4. Validates the returned ID to ensure it is non-empty.
//
// Parameters:
//   - ctx: the context for the database operation.
//   - skill: pointer to domain.Skill to persist. This method mutates
//     skill.CreatedAt and skill.UpdatedAt.
func (r *skillRepository) Create(ctx context.Context, skill *domain.Skill) (string, error) {
	if skill == nil {
		return "", errors.New("failed to validate skill: payload is nil")
	}

	if err := skill.ValidatePayload(); err != nil {
		return "", fmt.Errorf("failed to validate skill: %w", err)
	}

	id := utils.GenerateKey()
	now := r.timeProvider()

	skill.CreatedAt = now
	skill.UpdatedAt = now

	tableIdent := pgx.Identifier{r.skillTable}.Sanitize()
	query := fmt.Sprintf(
		`INSERT INTO %s
		(id, icon, hex_color, label, category, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		tableIdent,
	)

	var returnedID string
	err := r.databaseAPI.QueryRow(
		ctx,
		query,
		id,
		skill.Icon,
		skill.HexColor,
		skill.Label,
		skill.Category,
		skill.CreatedAt,
		skill.UpdatedAt,
	).Scan(&returnedID)

	if err != nil {
		return "", fmt.Errorf("failed to create skill: %w", err)
	}

	if returnedID == "" {
		return "", errors.New("invalid skill returned: ID missing")
	}

	return returnedID, nil
}

// Get retrieves the Skill with the given id from the repository.
// It returns a pointer to domain.Skill on success.
// If the provided id is empty, Get returns an error indicating a missing ID.
// If no row matches the id, the underlying pgx.ErrNoRows is wrapped and returned.
// After scanning the database row, the returned skill is validated via
// skill.ValidateResponse(); an error is returned if validation fails.
// The query selects id, icon, hex_color, label, category, created_at and updated_at
// from the repository's skill table.
func (r *skillRepository) Get(ctx context.Context, id string) (*domain.Skill, error) {
	if id == "" {
		return nil, fmt.Errorf("failed to get skill: ID missing")
	}

	var skill domain.Skill

	tableIdent := pgx.Identifier{r.skillTable}.Sanitize()
	query := fmt.Sprintf(
		`SELECT id, icon, hex_color, label, category, created_at, updated_at
		FROM %s
		WHERE id = $1`,
		tableIdent,
	)

	err := r.databaseAPI.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&skill.Id,
		&skill.Icon,
		&skill.HexColor,
		&skill.Label,
		&skill.Category,
		&skill.CreatedAt,
		&skill.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get skill: %w", err)
		}
		return nil, fmt.Errorf("failed to scan skill: %w", err)
	}

	if err := skill.ValidateResponse(); err != nil {
		return nil, fmt.Errorf("invalid skill returned: %w", err)
	}

	return &skill, nil
}

// Update updates an existing Skill in the repository. It validates the provided
// skill (ensuring the payload is non-nil, the Id is present, and ValidatePayload
// succeeds), sets the UpdatedAt timestamp using the repository's timeProvider,
// and performs an SQL UPDATE of the icon, hex_color, label, category and
// updated_at columns. The updated row is returned as a domain.Skill populated
// from the database (including created_at and updated_at). If no row matches
// the given Id, (nil, nil) is returned to indicate "not found". Any validation
// or database error is returned wrapped. The provided context is used for the
// database operation.
func (r *skillRepository) Update(ctx context.Context, skill *domain.Skill) (*domain.Skill, error) {
	if skill == nil {
		return nil, errors.New("failed to validate skill: payload is nil")
	}
	if skill.Id == "" {
		return nil, fmt.Errorf("failed to update skill: ID missing")
	}

	if err := skill.ValidatePayload(); err != nil {
		return nil, fmt.Errorf("failed to validate skill: %w", err)
	}

	updatedAt := r.timeProvider()

	var updatedSkill domain.Skill

	tableIdent := pgx.Identifier{r.skillTable}.Sanitize()
	query := fmt.Sprintf(
		`UPDATE %s
		SET icon=$2,
			hex_color=$3,
			label=$4,
			category=$5,
			updated_at=$6
		WHERE id=$1
		RETURNING id, icon, hex_color, label, category, created_at, updated_at`,
		tableIdent,
	)

	err := r.databaseAPI.QueryRow(
		ctx,
		query,
		skill.Id,
		skill.Icon,
		skill.HexColor,
		skill.Label,
		skill.Category,
		updatedAt,
	).Scan(
		&updatedSkill.Id,
		&updatedSkill.Icon,
		&updatedSkill.HexColor,
		&updatedSkill.Label,
		&updatedSkill.Category,
		&updatedSkill.CreatedAt,
		&updatedSkill.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update skill: %w", err)
	}

	if err := updatedSkill.ValidateResponse(); err != nil {
		return nil, fmt.Errorf("invalid skill returned: %w", err)
	}

	return &updatedSkill, nil
}

// Delete removes the skill with the given id from the repository's skill table.
// It requires a non-empty id and uses the provided context for cancellation and timeouts.
// If id is empty, Delete returns an error indicating the missing ID.
// The method executes a DELETE statement via the repository's database API and
// wraps any execution error. If no rows are affected by the DELETE, it returns
// pgx.ErrNoRows to indicate that no matching record was found.
func (r *skillRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("failed to delete skill: ID missing")
	}

	tableIdent := pgx.Identifier{r.skillTable}.Sanitize()
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", tableIdent)

	cmdTag, err := r.databaseAPI.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}
	if cmdTag == nil {
		return fmt.Errorf("failed to delete skill: nil command tag")
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// List retrieves a slice of domain.Skill from the repository using the provided filter.
//
// Behavior and defaults:
//   - If filter.Page <= 0 it defaults to 1.
//   - If filter.PageSize <= 0 or > 20 it defaults to 20 (maximum page size is 20).
//   - If filter.SortBy is nil it defaults to domain.CreatedAt.
//   - Sort direction is ascending by default; set filter.SortAscending = false for DESC.
//   - If filter.Category is non-nil, results are filtered by the given category.
//
// Query details:
//   - The query selects columns: id, icon, hex_color, label, category, created_at, updated_at
//     from the repository's skill table.
//   - Category filtering is applied via a parameterized WHERE clause (uses $1, $2, ... placeholders).
//   - ORDER BY uses the value of filter.SortBy verbatim as the column to sort on; callers must
//     ensure it is a valid column name to avoid SQL errors.
//   - LIMIT and OFFSET are applied for pagination (OFFSET = (Page-1) * PageSize).
//
// Execution and errors:
//   - The method executes the constructed query using r.databaseAPI.Query with the accumulated args.
//   - Rows are scanned into domain.Skill values and returned as a slice.
//   - On query, scan, or row iteration failures the method returns a wrapped error with context
//     ("failed to list skills", "failed to scan skill", or "row iteration error").
//
// Returns:
//   - ([]domain.Skill, nil) on success.
//   - (nil, error) on failure.
func (r *skillRepository) List(ctx context.Context, filter domain.SkillFilter) ([]domain.Skill, error) {
	// Set defaults if not provided
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 20 {
		filter.PageSize = 20
	}
	if filter.SortBy == nil {
		defaultSort := domain.CreatedAt
		filter.SortBy = &defaultSort
	}

	baseQuery := fmt.Sprintf(
		`SELECT id, icon, hex_color, label, category, created_at, updated_at FROM %s`,
		r.skillTable,
	)
	var conditions []string
	var args []any
	argIdx := 1

	// Add optional category filter
	if filter.Category != nil {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, *filter.Category)
		argIdx++
	}

	// Append WHERE clause if any filters exist
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add sorting
	sortOrder := "ASC"
	if !filter.SortAscending {
		sortOrder = "DESC"
	}
	baseQuery += fmt.Sprintf(" ORDER BY %s %s", *filter.SortBy, sortOrder)

	// Add pagination
	offset := (filter.Page - 1) * filter.PageSize
	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.PageSize, offset)

	// Execute query
	rows, err := r.databaseAPI.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list skills: %w", err)
	}
	defer rows.Close()

	var skills []domain.Skill
	for rows.Next() {
		var skill domain.Skill

		err := rows.Scan(
			&skill.Id,
			&skill.Icon,
			&skill.HexColor,
			&skill.Label,
			&skill.Category,
			&skill.CreatedAt,
			&skill.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}

		skills = append(skills, skill)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return skills, nil
}
