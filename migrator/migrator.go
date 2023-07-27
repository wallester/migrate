package migrator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/wallester/migrate/direction"

	"github.com/juju/errors"
	"github.com/wallester/migrate/driver"
	"github.com/wallester/migrate/file"
	"github.com/wallester/migrate/printer"
	"github.com/wallester/migrate/version"
)

// IMigrator represents possible migration actions
type IMigrator interface {
	Migrate(args Args) error
	Create(name, path string, verbose bool) (*file.Pair, error)
}

type Migrator struct {
	db     driver.IDriver
	output printer.IPrinter
}

var _ IMigrator = (*Migrator)(nil)

// New returns new instance
func New(db driver.IDriver, output printer.IPrinter) *Migrator {
	return &Migrator{
		db:     db,
		output: output,
	}
}

// Migrate migrates up or down
func (m *Migrator) Migrate(args Args) error {
	started := time.Now()

	files, err := file.ListFiles(args.Path, args.Direction)
	if err != nil {
		return errors.Annotate(err, "listing migration files failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), args.DBConnectionTimeoutDuration)
	defer cancel()
	if err := m.db.Open(ctx, args.URL); err != nil {
		return errors.Annotate(err, "opening database connection failed")
	}

	defer func() {
		if err := m.db.Close(); err != nil {
			m.output.Println(errors.Annotate(err, "closing database connection failed").Error())
		}
	}()

	migratedFiles, err := m.applyMigrations(files, args)
	if err != nil {
		return errors.Annotate(err, "migrating failed")
	}

	for _, f := range migratedFiles {
		m.output.Println(args.Direction.ToANSIColoredPrefix(), f.Base)
	}

	if args.Verbose {
		spent := time.Since(started).Seconds()
		m.output.Println(fmt.Sprintf("\n%.4f", spent), "seconds")
	}

	return nil
}

func (m *Migrator) Create(name, path string, verbose bool) (*file.Pair, error) {
	name = strings.ReplaceAll(name, " ", "_")
	v := time.Now().Unix()

	up := file.File{
		Version: v,
		Base:    fmt.Sprintf("%d_%s.%s.sql", v, name, direction.Up.ToString()),
	}
	if err := up.Create(path); err != nil {
		return nil, errors.Annotate(err, "writing up migration file failed")
	}

	down := file.File{
		Version: v,
		Base:    fmt.Sprintf("%d_%s.%s.sql", v, name, direction.Down.ToString()),
	}
	if err := down.Create(path); err != nil {
		return nil, errors.Annotate(err, "writing down migration file failed")
	}

	if verbose {
		m.output.Println("Version", v, "migration files created in", path)
		m.output.Println(up.Base)
		m.output.Println(down.Base)
	}

	return &file.Pair{
		Up:   up,
		Down: down,
	}, nil
}

// private

func (m *Migrator) applyMigrations(files []file.File, args Args) ([]file.File, error) {
	timeoutDuration := time.Duration(args.TimeoutSeconds * 1000)
	if timeoutDuration == 0 {
		timeoutDuration = args.TimeoutDuration
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	if err := m.db.CreateMigrationsTable(ctx); err != nil {
		return nil, errors.Annotate(err, "creating migrations table failed")
	}

	alreadyMigrated, err := m.db.SelectAllMigrations(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "selecting existing migrations failed")
	}

	needsMigration, err := m.chooseMigrations(files, alreadyMigrated, args)
	if err != nil {
		return nil, errors.Annotate(err, "choosing migrations failed")
	}

	if len(needsMigration) == 0 {
		if args.Verbose {
			m.output.Println("nothing to migrate")
		}

		return nil, nil
	}

	for _, f := range needsMigration {
		if err := m.db.Migrate(ctx, f, args.Direction); err != nil {
			return nil, errors.Annotatef(err, "applying migration failed: %s", f.Base)
		}
	}

	return needsMigration, nil
}

func (m *Migrator) chooseMigrations(files []file.File, alreadyMigrated version.Versions, args Args) ([]file.File, error) {
	maxMigratedVersion := alreadyMigrated.Max()
	boolDirection := bool(args.Direction)

	needsMigration := make([]file.File, 0, len(files))
	for _, f := range files {
		_, isMigrated := alreadyMigrated[f.Version]

		if boolDirection && isMigrated {
			continue
		}

		if !boolDirection && !isMigrated {
			continue
		}

		if boolDirection && maxMigratedVersion > f.Version && !args.NoVerify {
			return nil, fmt.Errorf("cannot migrate up %s, because it's older than already migrated version %d", f.Base, maxMigratedVersion)
		}

		needsMigration = append(needsMigration, f)
		if args.Verbose {
			m.output.Println(fmt.Sprintf("file %s needs migration", f.Base))
		}
	}

	totalFilesCount := len(needsMigration)
	if args.Verbose {
		m.output.Println(fmt.Sprintf("total files number to be migrated is %d", totalFilesCount))
	}

	if args.Steps > 0 && totalFilesCount >= args.Steps {
		if args.Verbose {
			m.output.Println(fmt.Sprintf("only %d files will be migrated out of total number", args.Steps))
		}

		needsMigration = needsMigration[:args.Steps]
	}

	return needsMigration, nil
}
