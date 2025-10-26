package v1

import (
	"context"
	"errors"
	"fmt"
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

	query := fmt.Sprintf(
		`INSERT INTO %s
		(id, icon, hex_color, label, category, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		r.skillTable,
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

	query := fmt.Sprintf(
		`SELECT id, icon, hex_color, label, category, created_at, updated_at
		FROM %s
		WHERE id = $1`,
		r.skillTable,
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

	now := r.timeProvider()
	skill.UpdatedAt = now

	var updatedSkill domain.Skill

	query := fmt.Sprintf(
		`UPDATE %s
		SET icon=$2,
			hex_color=$3,
			label=$4,
			category=$5,
			updated_at=$6
		WHERE id=$1
		RETURNING id, icon, hex_color, label, category, created_at, updated_at`,
		r.skillTable,
	)

	err := r.databaseAPI.QueryRow(
		ctx,
		query,
		skill.Id,
		skill.Icon,
		skill.HexColor,
		skill.Label,
		skill.Category,
		skill.UpdatedAt,
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

	return &updatedSkill, nil
}
