package migrator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/mgutz/ansi"
	"github.com/wallester/migrate/direction"
	"github.com/wallester/migrate/driver"
	"github.com/wallester/migrate/file"
	"github.com/wallester/migrate/printer"
	"github.com/wallester/migrate/version"
)

// IMigrator represents possible migration actions
type IMigrator interface {
	Migrate(args Args) error
	Create(name, path string) (*file.Pair, error)
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

	files, err := file.ListFiles(args.Path, args.Up)
	if err != nil {
		return errors.Annotate(err, "listing migration files failed")
	}

	if err := m.db.Open(args.URL); err != nil {
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
		m.output.Println(printPrefix[args.Up], f.Base)
	}

	m.output.Println("")
	spent := time.Since(started).Seconds()
	m.output.Println(fmt.Sprintf("%.4f", spent), "seconds")

	return nil
}

func (m *Migrator) Create(name, path string) (*file.Pair, error) {
	name = strings.ReplaceAll(name, " ", "_")
	v := time.Now().Unix()

	up := file.File{
		Version: v,
		Base:    fmt.Sprintf("%d_%s.up.sql", v, name),
	}
	if err := up.Create(path); err != nil {
		return nil, errors.Annotate(err, "writing up migration file failed")
	}

	down := file.File{
		Version: v,
		Base:    fmt.Sprintf("%d_%s.down.sql", v, name),
	}
	if err := down.Create(path); err != nil {
		return nil, errors.Annotate(err, "writing down migration file failed")
	}

	m.output.Println("Version", v, "migration files created in", path)
	m.output.Println(up.Base)
	m.output.Println(down.Base)

	return &file.Pair{
		Up:   up,
		Down: down,
	}, nil
}

// private

var printPrefix = map[bool]string{
	direction.Up:   ansi.Green + ">" + ansi.Reset,
	direction.Down: ansi.Red + "<" + ansi.Reset,
}

func (m *Migrator) applyMigrations(files []file.File, args Args) ([]file.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(args.TimeoutSeconds)*time.Second)
	defer cancel()

	if err := m.db.CreateMigrationsTable(ctx); err != nil {
		return nil, errors.Annotate(err, "creating migrations table failed")
	}

	alreadyMigrated, err := m.db.SelectAllMigrations(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "selecting existing migrations failed")
	}

	needsMigration, err := chooseMigrations(files, alreadyMigrated, args)
	if err != nil {
		return nil, errors.Annotate(err, "choosing migrations failed")
	}

	if len(needsMigration) == 0 {
		return nil, nil
	}

	for _, f := range needsMigration {
		if err := m.db.Migrate(ctx, f, args.Up); err != nil {
			return nil, errors.Annotatef(err, "applying migration failed: %s", f.Base)
		}
	}

	return needsMigration, nil
}

func chooseMigrations(files []file.File, alreadyMigrated version.Versions, args Args) ([]file.File, error) {
	maxMigratedVersion := alreadyMigrated.Max()

	needsMigration := make([]file.File, 0, len(files))
	for _, f := range files {
		_, isMigrated := alreadyMigrated[f.Version]

		if args.Up && isMigrated {
			continue
		}

		if !args.Up && !isMigrated {
			continue
		}

		if args.Up && maxMigratedVersion > f.Version && !args.NoVerify {
			return nil, fmt.Errorf("cannot migrate up %s, because it's older than already migrated version %d", f.Base, maxMigratedVersion)
		}

		needsMigration = append(needsMigration, f)
	}

	if args.Steps > 0 && len(needsMigration) >= args.Steps {
		needsMigration = needsMigration[:args.Steps]
	}

	return needsMigration, nil
}
