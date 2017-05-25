package postgres

import (
	"context"
	"database/sql"

	"github.com/juju/errors"
	_ "github.com/lib/pq" // import driver
	"github.com/wallester/migrate/direction"
	"github.com/wallester/migrate/driver"
	"github.com/wallester/migrate/file"
)

type postgres struct {
	connection *sql.DB
}

// New returns new instance
func New() driver.Driver {
	return &postgres{}
}

// Open opens database connection
func (db *postgres) Open(url string) error {
	connection, err := sql.Open("postgres", url)
	if err != nil {
		return errors.Annotate(err, "connecting to database failed")
	}

	db.connection = connection

	return nil
}

// Close closes database connection
func (db *postgres) Close() error {
	if err := db.connection.Close(); err != nil {
		return errors.Annotate(err, "closing database connection failed")
	}

	return nil
}

// SelectAllMigrations selects existing migrations
func (db *postgres) SelectAllMigrations(ctx context.Context) (map[int64]bool, error) {
	rows, err := db.connection.QueryContext(ctx, `
		SELECT version FROM schema_migrations
	`)
	if err != nil {
		return nil, errors.Annotate(err, "selecting existing migration versions failed")
	}

	migrated := make(map[int64]bool)
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			if closeErr := rows.Close(); closeErr != nil {
				return nil, errors.Annotate(closeErr, "closing rows failed")
			}

			return nil, errors.Annotate(err, "scanning version failed")
		}

		migrated[version] = true
	}

	if closeErr := rows.Close(); closeErr != nil {
		return nil, errors.Annotate(closeErr, "closing rows failed")
	}

	return migrated, nil
}

// CreateMigrationsTable creates migrations table if it does not exist yet
func (db *postgres) CreateMigrationsTable(ctx context.Context) error {
	if _, err := db.connection.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations(
			version bigint not null primary key
		)
	`); err != nil {
		return errors.Annotate(err, "creating schema_migrations table failed")
	}

	return nil
}

var applyMigrationSQL = map[bool]string{
	direction.Up:   "INSERT INTO schema_migrations(version) VALUES($1)",
	direction.Down: "DELETE FROM schema_migrations WHERE version = $1",
}

// ApplyMigrations applies migrations to database
func (db *postgres) ApplyMigrations(ctx context.Context, files []file.File, up bool) error {
	tx, err := db.connection.Begin()
	if err != nil {
		return errors.Annotate(err, "starting database transaction failed")
	}

	for _, file := range files {
		if _, err := tx.ExecContext(ctx, file.SQL); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Annotate(err, "rolling back transaction failed")
			}

			return errors.Annotate(err, "executing migration failed")
		}

		if _, err := tx.ExecContext(ctx, applyMigrationSQL[up], file.Version); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Annotate(err, "rolling back transaction failed")
			}

			return errors.Annotate(err, "executing migration failed")
		}
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Annotate(err, "rolling back transaction failed")
		}

		return errors.Annotate(err, "committing migrations failed")
	}

	return nil
}
