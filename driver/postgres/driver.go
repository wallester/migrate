package postgres

import (
	"context"
	"database/sql"

	"github.com/juju/errors"
	_ "github.com/lib/pq" // import driver
	"github.com/wallester/migrate/direction"
	"github.com/wallester/migrate/driver"
	"github.com/wallester/migrate/file"
	"github.com/wallester/migrate/version"
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
func (db *postgres) SelectAllMigrations(ctx context.Context) (version.Versions, error) {
	rows, err := db.connection.QueryContext(ctx, `
		SELECT version FROM schema_migrations
	`)
	if err != nil {
		return nil, errors.Annotate(err, "selecting existing migration versions failed")
	}

	var exists struct{}
	migrated := make(version.Versions)
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			if closeErr := rows.Close(); closeErr != nil {
				return nil, errors.Annotate(closeErr, "closing rows failed")
			}

			return nil, errors.Annotate(err, "scanning version failed")
		}

		migrated[version] = exists
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
			version bigint not null primary key,
			applied_at timestamp without time zone
		)
	`); err != nil {
		return errors.Annotate(err, "creating schema_migrations table failed")
	}

	var appliedAtExists bool
	if err := db.connection.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT
				1
			FROM
				information_schema.columns
			WHERE
				table_name = 'schema_migrations'
			AND
				column_name = 'applied_at'
		)
	`).Scan(&appliedAtExists); err != nil {
		return errors.Annotate(err, "checking if applied_at timestamp exists failed")
	}

	if !appliedAtExists {
		if _, err := db.connection.ExecContext(ctx, `
			ALTER TABLE schema_migrations ADD COLUMN applied_at timestamp without time zone
		`); err != nil {
			return errors.Annotate(err, "adding applied_at timestamp failed")
		}
	}

	return nil
}

var applyMigrationSQL = map[bool]string{
	direction.Up:   "INSERT INTO schema_migrations(version, applied_at) VALUES($1, NOW() at time zone 'utc')",
	direction.Down: "DELETE FROM schema_migrations WHERE version = $1",
}

// ApplyMigrations applies migrations to database
func (db *postgres) ApplyMigrations(ctx context.Context, files []file.File, up bool) error {
	tx, err := db.connection.Begin()
	if err != nil {
		return errors.Annotate(err, "starting database transaction failed")
	}

	rollback := func(reasonErr error) error {
		if err := tx.Rollback(); err != nil {
			return errors.Annotate(err, "rolling back transaction failed")
		}

		return reasonErr
	}

	for _, file := range files {
		if _, err := tx.ExecContext(ctx, file.SQL); err != nil {
			return rollback(errors.Annotatef(err, "executing %s migration failed", file.Base))
		}

		if _, err := tx.ExecContext(ctx, applyMigrationSQL[up], file.Version); err != nil {
			return rollback(errors.Annotatef(err, "executing %s migration failed", file.Base))
		}
	}

	if err := tx.Commit(); err != nil {
		return rollback(errors.Annotate(err, "committing migrations failed"))
	}

	return nil
}
