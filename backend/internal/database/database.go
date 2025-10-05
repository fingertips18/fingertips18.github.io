package database

import (
	"context"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/client"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// --- Wrappers ---

type rowWrapper struct {
	pgx.Row
}

func (r rowWrapper) Scan(dest ...any) error {
	return r.Row.Scan(dest...)
}

type rowsWrapper struct {
	pgx.Rows
}

func (r rowsWrapper) Next() bool             { return r.Rows.Next() }
func (r rowsWrapper) Scan(dest ...any) error { return r.Rows.Scan(dest...) }
func (r rowsWrapper) Close()                 { r.Rows.Close() }
func (r rowsWrapper) Err() error             { return r.Rows.Err() }

type commandTagWrapper struct {
	pgconn.CommandTag
}

func (c commandTagWrapper) RowsAffected() int64 {
	return c.CommandTag.RowsAffected()
}

// --- Interfaces ---

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}

type CommandTag interface {
	RowsAffected() int64
}

type DatabaseAPI interface {
	QueryRow(ctx context.Context, query string, args ...any) Row
	Exec(ctx context.Context, query string, args ...any) (CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	Close()
}

type database struct {
	pool client.PgxAPI
}

// NewDatabase creates a new instance of Database using the provided connection string.
// It initializes a connection pool to the PostgreSQL database and returns a pointer to the Database struct.
//
// Parameters:
//   - connectionString: the database connection string.
//
// Returns:
//   - *Database: a pointer to the initialized Database instance.
func NewDatabase(connectionString string) DatabaseAPI {
	pool := client.NewPGXAPI(connectionString)

	return &database{
		pool: pool,
	}
}

func (d *database) QueryRow(ctx context.Context, query string, args ...any) Row {
	return rowWrapper{d.pool.QueryRow(ctx, query, args...)}
}

func (d *database) Exec(ctx context.Context, query string, args ...any) (CommandTag, error) {
	tag, err := d.pool.Exec(ctx, query, args...)
	return commandTagWrapper{tag}, err
}

func (d *database) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	rows, err := d.pool.Query(ctx, query, args...)
	return rowsWrapper{rows}, err
}

func (d *database) Close() {
	d.pool.Close()
}
