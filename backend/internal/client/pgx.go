package client

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxAPI interface {
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	Close()
}

type pgxClient struct {
	pool *pgxpool.Pool
}

// NewPGXAPI creates a new instance of PGXAPI using the provided PostgreSQL connection string.
// It establishes a connection pool to the database and returns a PGXAPI implementation.
// If the connection cannot be established, the function logs a fatal error and terminates the application.
//
// Parameters:
//   - connectionString: The PostgreSQL connection string.
//
// Returns:
//   - PGXAPI: An implementation of the PGXAPI interface connected to the specified database.
func NewPGXAPI(connectionString string) PgxAPI {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connectionString)

	if err != nil {
		log.Fatalf("Unable to connect database: %v\n", err)
	}

	return &pgxClient{
		pool: pool,
	}
}

// QueryRow executes a query that is expected to return at most one row.
// It takes a context, a SQL query string, and any number of arguments for the query placeholders.
// It returns a pgx.Row, which can be used to scan the result.
// If no rows are returned, pgx.Row's Scan will return pgx.ErrNoRows.
func (p *pgxClient) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, query, args...)
}

// Exec executes the given SQL query with the provided arguments using the connection pool.
// It returns a pgconn.CommandTag, which contains information about the command executed,
// and an error if the execution fails.
//
// Parameters:
//   - ctx: The context for controlling cancellation and deadlines.
//   - query: The SQL query string to execute.
//   - args: Variadic arguments to be substituted into the query.
//
// Returns:
//   - pgconn.CommandTag: Information about the executed command.
//   - error: An error if the execution fails, otherwise nil.
func (p *pgxClient) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, query, args...)
}

// Query executes a query against the PostgreSQL database using the provided context,
// SQL query string, and optional arguments. It returns the resulting rows and any error encountered.
// The method delegates the query execution to the underlying connection pool.
func (p *pgxClient) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, query, args...)
}

// Close releases all resources used by the pgxClient by closing its underlying connection pool.
func (p *pgxClient) Close() {
	p.pool.Close()
}
