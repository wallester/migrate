package driver

import (
	"context"
	"database/sql"
	"log"

	"github.com/juju/errors"
	_ "github.com/lib/pq" // import driver
	"github.com/wallester/migrate/file"
)

// Driver represents database driver
type Driver interface {
	Open(url string) error
	CreateMigrationsTable(ctx context.Context) error
	SelectMigrations(ctx context.Context) (map[int]bool, error)
	ApplyMigrations(ctx context.Context, files []file.File, up bool) error
	Close()
}

type driver struct {
	connection *sql.DB
}

// New returns new instance
func New() Driver {
	return &driver{}
}

// Open opens database connection
func (db *driver) Open(url string) error {
	connection, err := sql.Open("postgres", url)
	if err != nil {
		return errors.Annotate(err, "connecting to database failed")
	}

	db.connection = connection

	return nil
}

// Close closes database connection
func (db *driver) Close() {
	if err := db.connection.Close(); err != nil {
		log.Println("Warning", errors.Annotate(err, "closing database connection failed"))
	}
}

func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Println("Warning", errors.Annotate(err, "closing sql result rows failed"))
	}
}

// SelectMigrations selects existing migrations
func (db *driver) SelectMigrations(ctx context.Context) (map[int]bool, error) {
	rows, err := db.connection.QueryContext(ctx, `
		SELECT version FROM schema_migrations
	`)
	if err != nil {
		return nil, errors.Annotate(err, "selecting existing migration versions failed")
	}

	defer closeRows(rows)

	migrated := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, errors.Annotate(err, "scanning version failed")
		}

		migrated[version] = true
	}

	return migrated, nil
}

// CreateMigrationsTable creates migrations table if it does not exist yet
func (db *driver) CreateMigrationsTable(ctx context.Context) error {
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
	true:  "INSERT INTO schema_migrations(version) VALUES($1)",
	false: "DELETE FROM schema_migrations WHERE version = $1",
}

// ApplyMigrations applies migrations to database
func (db *driver) ApplyMigrations(ctx context.Context, files []file.File, up bool) error {
	tx, err := db.connection.Begin()
	if err != nil {
		return errors.Annotate(err, "starting database transaction failed")
	}

	for _, file := range files {
		if _, err := tx.ExecContext(ctx, file.SQL); err != nil {
			return errors.Annotate(err, "executing migration failed")
		}

		if _, err := tx.ExecContext(ctx, applyMigrationSQL[up], file.Version); err != nil {
			return errors.Annotate(err, "executing migration failed")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Annotate(err, "committing migrations failed")
	}

	return nil
}
