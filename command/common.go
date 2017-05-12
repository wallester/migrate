package command

import (
	"context"
	"fmt"
	"time"

	"github.com/juju/errors"
	"github.com/mgutz/ansi"
	"github.com/urfave/cli"
	"github.com/wallester/migrate/file"
	"github.com/wallester/migrate/flag"
)

const timeoutSeconds = 1

var printPrefix = map[bool]string{
	true:  ansi.Green + ">" + ansi.Reset,
	false: ansi.Red + "<" + ansi.Reset,
}

func migrate(c *cli.Context, up bool) error {
	started := time.Now()

	path := flag.Get(c, flag.FlagPath)
	if path == "" {
		return flag.NewRequiredFlagError(flag.FlagPath)
	}

	files, err := file.ListFiles(path, up)
	if err != nil {
		return errors.Annotate(err, "listing migration files failed")
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
		fmt.Println(printPrefix[up], file.Base)
	}

	fmt.Println("")
	spent := time.Since(started).Seconds()
	fmt.Println(fmt.Sprintf("%.4f", spent), "seconds")

	return nil
}

func migrateFiles(url string, files []file.File, up bool) ([]file.File, error) {
	db, err := openDB(url)
	if err != nil {
		return nil, errors.Annotate(err, "opening database connection failed")
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

	if err := applyMigrations(ctx, db, needsMigration, up); err != nil {
		return nil, errors.Annotate(err, "applying migrations failed")
	}

	return needsMigration, nil
}

func chooseMigrations(files []file.File, alreadyMigrated map[int]bool, up bool) ([]file.File, error) {
	var needsMigration []file.File
	for _, file := range files {
		if (up && !alreadyMigrated[file.Version]) || (!up && alreadyMigrated[file.Version]) {
			needsMigration = append(needsMigration, file)
		}
	}

	return needsMigration, nil
}
