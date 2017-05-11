package command

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/juju/errors"
	_ "github.com/lib/pq" // import driver
	"github.com/mgutz/ansi"
	"github.com/urfave/cli"
	"github.com/wallester/migrate/flag"
)

const timeoutSeconds = 1

var (
	filePrefix = map[bool]string{
		true:  "up",
		false: "down",
	}
	printPrefix = map[bool]string{
		true:  ansi.Green + ">" + ansi.Reset,
		false: ansi.Red + "<" + ansi.Reset,
	}
	applyMigrationSQL = map[bool]string{
		true:  "INSERT INTO schema_migrations(version) VALUES($1)",
		false: "DELETE FROM schema_migrations WHERE version = $1",
	}
)

func migrate(c *cli.Context, up bool) error {
	started := time.Now()

	path := flag.Get(c, flag.FlagPath)
	if path == "" {
		return flag.NewRequiredFlagError(flag.FlagPath)
	}

	files, err := filepath.Glob(filepath.Join(path, "*_*."+filePrefix[up]+".sql"))
	if err != nil {
		return errors.Annotate(err, "getting migration files failed")
	}

	if up {
		sort.Sort(sort.StringSlice(files))
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(files)))
	}

	url := flag.Get(c, flag.FlagURL)
	if url == "" {
		return flag.NewRequiredFlagError(flag.FlagURL)
	}

	migratedFiles, err := migrateFiles(url, files, up)
	if err != nil {
		return errors.Annotate(err, "migrating failed")
	}

	for _, file := range migratedFiles {
		fmt.Println(printPrefix[up], filepath.Base(file))
	}

	fmt.Println("")
	spent := time.Since(started).Seconds()
	fmt.Println(fmt.Sprintf("%.4f", spent), "seconds")

	return nil
}

func migrateFiles(url string, files []string, up bool) ([]string, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errors.Annotate(err, "connecting to database failed")
	}

	defer closeDB(db)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancel()

	err = createMigrationsTable(ctx, db)
	if err != nil {
		return nil, errors.Annotate(err, "creating migrations table failed")
	}

	alreadyMigrated, err := selectExistingMigrations(ctx, db)
	if err != nil {
		return nil, errors.Annotate(err, "selecting existing migrations failed")
	}

	needsMigration, err := chooseMigrations(files, alreadyMigrated, up)
	if err != nil {
		return nil, errors.Annotate(err, "choosing migrations failed")
	}

	if err := applyMigrations(ctx, needsMigration, db, up); err != nil {
		return nil, errors.Annotate(err, "applying migrations failed")
	}

	return needsMigration, nil
}

func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Println("Warning", errors.Annotate(err, "closing database connection failed"))
	}
}

func parseVersion(file string) (int, error) {
	return strconv.Atoi(strings.Split(filepath.Base(file), "_")[0])
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

func migrateFile(ctx context.Context, file string, tx *sql.Tx, up bool) error {
	version, err := parseVersion(file)
	if err != nil {
		return errors.Annotate(err, "parsing file version failed")
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Annotate(err, "reading migration file failed")
	}

	if _, err := tx.ExecContext(ctx, string(b)); err != nil {
		return errors.Annotate(err, "executing migration failed")
	}

	if _, err := tx.ExecContext(ctx, applyMigrationSQL[up], version); err != nil {
		return errors.Annotate(err, "executing migration failed")
	}

	return nil
}

func chooseMigrations(files []string, alreadyMigrated map[int]bool, up bool) ([]string, error) {
	var needsMigration []string
	for _, file := range files {
		version, err := parseVersion(file)
		if err != nil {
			return nil, errors.Annotate(err, "parsing file version failed")
		}

		if (up && !alreadyMigrated[version]) || (!up && alreadyMigrated[version]) {
			needsMigration = append(needsMigration, file)
		}
	}

	return needsMigration, nil
}

func applyMigrations(ctx context.Context, needsMigration []string, db *sql.DB, up bool) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.Annotate(err, "starting database transaction failed")
	}

	for _, file := range needsMigration {
		if err := migrateFile(ctx, file, tx, up); err != nil {
			return errors.Annotate(err, "migrating file failed")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Annotate(err, "committing migrations failed")
	}

	return nil
}

func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Println("Warning", errors.Annotate(err, "closing sql result rows failed"))
	}
}
