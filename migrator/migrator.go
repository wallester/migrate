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

// Migrator represents possible migration actions
type Migrator interface {
	Migrate(path string, url string, dir direction.Direction, steps int, timeoutSeconds int) error
	Create(name string, path string) (*file.Pair, error)
}

type migrator struct {
	db     driver.Driver
	output printer.Printer
}

// New returns new instance
func New(db driver.Driver, output printer.Printer) Migrator {
	return &migrator{
		db:     db,
		output: output,
	}
}

var printPrefix = map[direction.Direction]string{
	direction.Up:   ansi.Green + ">" + ansi.Reset,
	direction.Down: ansi.Red + "<" + ansi.Reset,
}

// Migrate migrates up or down
func (m *migrator) Migrate(path string, url string, dir direction.Direction, steps int, timeoutSeconds int) error {
	started := time.Now()

	files, err := file.ListFiles(path, dir)
	if err != nil {
		return errors.Annotate(err, "listing migration files failed")
	}

	err = m.db.Open(url)
	if err != nil {
		return errors.Annotate(err, "opening database connection failed")
	}

	migratedFiles, err := m.applyMigrations(files, dir, steps, timeoutSeconds)
	if err != nil {
		if closeErr := m.db.Close(); closeErr != nil {
			return errors.Annotate(closeErr, "closing database connection failed")
		}

		return errors.Annotate(err, "migrating failed")
	}

	for _, file := range migratedFiles {
		m.output.Println(printPrefix[dir], file.Base)
	}

	m.output.Println("")
	spent := time.Since(started).Seconds()
	m.output.Println(fmt.Sprintf("%.4f", spent), "seconds")

	if closeErr := m.db.Close(); closeErr != nil {
		return errors.Annotate(closeErr, "closing database connection failed")
	}

	return nil
}

func (m *migrator) applyMigrations(files []file.File, dir direction.Direction, steps int, timeoutSeconds int) ([]file.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	if err := m.db.CreateMigrationsTable(ctx); err != nil {
		return nil, errors.Annotate(err, "creating migrations table failed")
	}

	alreadyMigrated, err := m.db.SelectAllMigrations(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "selecting existing migrations failed")
	}

	needsMigration, err := chooseMigrations(files, alreadyMigrated, dir, steps)
	if err != nil {
		return nil, errors.Annotate(err, "choosing migrations failed")
	}

	if len(needsMigration) > 0 {
		if err := m.db.ApplyMigrations(ctx, needsMigration, dir); err != nil {
			return nil, errors.Annotate(err, "applying migrations failed")
		}
	}

	return needsMigration, nil
}

func chooseMigrations(files []file.File, alreadyMigrated version.Versions, dir direction.Direction, steps int) ([]file.File, error) {
	maxMigratedVersion := alreadyMigrated.Max()

	var needsMigration []file.File
	for _, f := range files {
		_, isMigrated := alreadyMigrated[f.Version]

		if dir == direction.Up && isMigrated {
			continue
		}

		if dir == direction.Down && !isMigrated {
			continue
		}

		if dir == direction.Up && maxMigratedVersion > f.Version {
			return nil, fmt.Errorf("cannot migrate up %s, because it's older than already migrated version %d", f.Base, maxMigratedVersion)
		}

		needsMigration = append(needsMigration, f)
	}

	if steps > 0 && len(needsMigration) >= steps {
		needsMigration = needsMigration[:steps]
	}

	return needsMigration, nil
}

func (m *migrator) Create(name string, path string) (*file.Pair, error) {
	name = strings.Replace(name, " ", "_", -1)
	version := time.Now().Unix()

	up := file.File{
		Version: version,
		Base:    fmt.Sprintf("%d_%s.up.sql", version, name),
	}
	if err := up.Create(path); err != nil {
		return nil, errors.Annotate(err, "writing up migration file failed")
	}

	down := file.File{
		Version: version,
		Base:    fmt.Sprintf("%d_%s.down.sql", version, name),
	}
	if err := down.Create(path); err != nil {
		return nil, errors.Annotate(err, "writing down migration file failed")
	}

	m.output.Println("Version", version, "migration files created in", path)
	m.output.Println(up.Base)
	m.output.Println(down.Base)

	return &file.Pair{
		Up:   up,
		Down: down,
	}, nil
}
