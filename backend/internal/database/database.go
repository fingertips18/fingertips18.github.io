package database

import (
	"github.com/fingertips18/fingertips18.github.io/backend/internal/client"
)

type Database struct {
	Pool client.PGXAPI
}

// NewDatabase creates a new instance of Database using the provided connection string.
// It initializes a connection pool to the PostgreSQL database and returns a pointer to the Database struct.
//
// Parameters:
//   - connectionString: the database connection string.
//
// Returns:
//   - *Database: a pointer to the initialized Database instance.
func NewDatabase(connectionString string) *Database {
	pool := client.NewPGXAPI(connectionString)

	return &Database{
		Pool: pool,
	}
}
