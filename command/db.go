package command

import (
	"context"
	"database/sql"
	"log"

	"github.com/juju/errors"
	_ "github.com/lib/pq" // import driver
)

var applyMigrationSQL = map[bool]string{
	true:  "INSERT INTO schema_migrations(version) VALUES($1)",
	false: "DELETE FROM schema_migrations WHERE version = $1",
}

func openDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errors.Annotate(err, "connecting to database failed")
	}

	return db, nil
}

func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Println("Warning", errors.Annotate(err, "closing database connection failed"))
	}
}

func selectExistingMigrations(ctx context.Context, db *sql.DB) (map[int]bool, error) {
	rows, err := db.QueryContext(ctx, `
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

func createMigrationsTable(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations(
			version bigint not null primary key
		)
	`); err != nil {
		return errors.Annotate(err, "creating schema_migrations table failed")
	}

	return nil
}

func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Println("Warning", errors.Annotate(err, "closing sql result rows failed"))
	}
}

func applyMigrations(ctx context.Context, db *sql.DB, files []MigrationFile, up bool) error {
	tx, err := db.Begin()
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