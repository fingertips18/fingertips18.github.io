package v1

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/utils"
)

type SkillRepository interface {
	Create(ctx context.Context, education *domain.CreateSkill) (string, error)
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

// Create inserts a new skill record into the repository and returns the generated ID.
//
// It performs the following steps:
//  1. Validates the provided CreateSkill payload via skill.ValidatePayload().
//  2. Generates a new ID using utils.GenerateKey().
//  3. Sets skill.CreatedAt and skill.UpdatedAt using the repository's timeProvider.
//  4. Executes an INSERT ... RETURNING id against the repository's skillTable and
//     scans the returned id.
//
// Parameters:
//   - ctx: context for cancellation and timeouts propagated to the database call.
//   - skill: pointer to domain.CreateSkill to persist. This method mutates
//     skill.CreatedAt and skill.UpdatedAt; callers should provide a non-nil payload.
//
// Returns:
//   - (string, error): the newly created skill ID on success, or a non-nil error on failure.
//     Error cases include payload validation failure, database query/scan errors, or an
//     unexpected empty ID returned by the database.
func (r *skillRepository) Create(ctx context.Context, skill *domain.CreateSkill) (string, error) {
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
