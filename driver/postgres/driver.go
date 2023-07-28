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

type Postgres struct {
	connection *sql.DB
}

var _ driver.IDriver = (*Postgres)(nil)

// New returns new instance
func New() *Postgres {
	return &Postgres{}
}

// Open opens database connection
func (db *Postgres) Open(ctx context.Context, url string) error {
	connection, err := sql.Open("postgres", url)
	if err != nil {
		return errors.Annotate(err, "connecting to database failed")
	}

	if err := connection.PingContext(ctx); err != nil {
		return errors.Annotate(err, "pinging database failed")
	}

	db.connection = connection

	return nil
}

// Close closes database connection
func (db *Postgres) Close() error {
	if err := db.connection.Close(); err != nil {
		return errors.Annotate(err, "closing database connection failed")
	}

	return nil
}

// SelectAllMigrations selects existing migrations
func (db *Postgres) SelectAllMigrations(ctx context.Context) (version.Versions, error) {
	rows, err := db.connection.QueryContext(ctx, `
		SELECT version FROM schema_migrations
	`)
	if err != nil {
		return nil, errors.Annotate(err, "selecting existing migration versions failed")
	}

	var exists struct{}
	migrated := make(version.Versions)
	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			if err := rows.Close(); err != nil {
				return nil, errors.Annotate(err, "closing rows failed")
			}

			return nil, errors.Annotate(err, "scanning version failed")
		}

		migrated[v] = exists
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	if err := rows.Close(); err != nil {
		return nil, errors.Annotate(err, "closing rows failed")
	}

	return migrated, nil
}

// CreateMigrationsTable creates migrations table if it does not exist yet
func (db *Postgres) CreateMigrationsTable(ctx context.Context) error {
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

func (db *Postgres) Migrate(ctx context.Context, f file.File, d direction.Direction) error {
	tx, err := db.connection.BeginTx(ctx, nil)
	if err != nil {
		return errors.Annotate(err, "starting database transaction failed")
	}

	rollback := func(reasonErr error) error {
		if err := tx.Rollback(); err != nil {
			return errors.Annotate(err, "rolling back transaction failed")
		}

		return reasonErr
	}

	if _, err := tx.ExecContext(ctx, f.SQL); err != nil {
		return rollback(errors.Annotatef(err, "executing %s migration failed", f.Base))
	}

	if _, err := tx.ExecContext(ctx, applyMigrationSQL[d], f.Version); err != nil {
		return rollback(errors.Annotatef(err, "executing %s migration failed", f.Base))
	}

	if err := tx.Commit(); err != nil {
		return rollback(errors.Annotate(err, "committing migrations failed"))
	}

	return nil
}

// private

var applyMigrationSQL = map[direction.Direction]string{
	direction.Up:   "INSERT INTO schema_migrations(version, applied_at) VALUES($1, NOW() at time zone 'utc')",
	direction.Down: "DELETE FROM schema_migrations WHERE version = $1",
}
