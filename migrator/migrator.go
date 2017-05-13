package migrator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/mgutz/ansi"
	"github.com/wallester/migrate/driver"
	"github.com/wallester/migrate/file"
	"github.com/wallester/migrate/printer"
)

// Migrator represents possible migration actions
type Migrator interface {
	MigrateAll(path string, url string, up bool) error
	Create(name string, path string) (*file.Pair, error)
}

type migrator struct {
	db driver.Driver
	p  printer.Printer
}

// New returns new instance
func New(db driver.Driver, p printer.Printer) Migrator {
	return &migrator{
		db: db,
		p:  p,
	}
}

var printPrefix = map[bool]string{
	true:  ansi.Green + ">" + ansi.Reset,
	false: ansi.Red + "<" + ansi.Reset,
}

// Migrate migrates all up or down
func (m *migrator) MigrateAll(path string, url string, up bool) error {
	started := time.Now()

	files, err := file.ListFiles(path, up)
	if err != nil {
		return errors.Annotate(err, "listing migration files failed")
	}

	err = m.db.Open(url)
	if err != nil {
		return errors.Annotate(err, "opening database connection failed")
	}

	defer m.db.Close()

	migratedFiles, err := m.executeFiles(files, up)
	if err != nil {
		return errors.Annotate(err, "migrating failed")
	}

	for _, file := range migratedFiles {
		m.p.Println(printPrefix[up], file.Base)
	}

	m.p.Println("")
	spent := time.Since(started).Seconds()
	m.p.Println(fmt.Sprintf("%.4f", spent), "seconds")

	return nil
}

const timeoutSeconds = 1

func (m *migrator) executeFiles(files []file.File, up bool) ([]file.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancel()

	if err := m.db.CreateMigrationsTable(ctx); err != nil {
		return nil, errors.Annotate(err, "creating migrations table failed")
	}

	alreadyMigrated, err := m.db.SelectMigrations(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "selecting existing migrations failed")
	}

	needsMigration, err := chooseMigrations(files, alreadyMigrated, up)
	if err != nil {
		return nil, errors.Annotate(err, "choosing migrations failed")
	}

	if len(needsMigration) > 0 {
		if err := m.db.ApplyMigrations(ctx, needsMigration, up); err != nil {
			return nil, errors.Annotate(err, "applying migrations failed")
		}
	}

	return needsMigration, nil
}

func chooseMigrations(files []file.File, alreadyMigrated map[int64]bool, up bool) ([]file.File, error) {
	var needsMigration []file.File
	for _, file := range files {
		if (up && !alreadyMigrated[file.Version]) || (!up && alreadyMigrated[file.Version]) {
			needsMigration = append(needsMigration, file)
		}
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

	m.p.Println("Version", version, "migration files created in", path)
	m.p.Println(up.Base)
	m.p.Println(down.Base)

	return &file.Pair{
		Up:   up,
		Down: down,
	}, nil
}
