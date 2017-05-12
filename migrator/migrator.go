package migrator

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/mgutz/ansi"
	"github.com/wallester/migrate/driver"
	"github.com/wallester/migrate/file"
)

// Migrator represents possible migration actions
type Migrator interface {
	Up(path string, url string) error
	Down(path string, url string) error
	Create(name string, path string) error
}

type migrator struct {
	db driver.Driver
}

// New returns new instance
func New(db driver.Driver) Migrator {
	return &migrator{db}
}

// Up migrates up
func (m *migrator) Up(path, url string) error {
	if err := m.execute(path, url, true); err != nil {
		return errors.Annotate(err, "migrating up failed")
	}

	return nil
}

// Down migrates down
func (m *migrator) Down(path, url string) error {
	if err := m.execute(path, url, false); err != nil {
		return errors.Annotate(err, "migrating down failed")
	}

	return nil
}

var printPrefix = map[bool]string{
	true:  ansi.Green + ">" + ansi.Reset,
	false: ansi.Red + "<" + ansi.Reset,
}

func (m *migrator) execute(path string, url string, up bool) error {
	started := time.Now()

	files, err := file.ListFiles(path, up)
	if err != nil {
		return errors.Annotate(err, "listing migration files failed")
	}

	if err := m.db.Open(url); err != nil {
		return errors.Annotate(err, "opening database connection failed")
	}

	defer m.db.Close()

	migratedFiles, err := m.executeFiles(files, up)
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

	if err := m.db.ApplyMigrations(ctx, needsMigration, up); err != nil {
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

func (m *migrator) Create(name string, path string) error {
	name = strings.Replace(name, " ", "_", -1)
	version := time.Now().Unix()

	up := fmt.Sprintf("%d_%s.up.sql", version, name)
	if err := ioutil.WriteFile(filepath.Join(path, up), nil, 0644); err != nil {
		return errors.Annotate(err, "writing up migration file failed")
	}

	down := fmt.Sprintf("%d_%s.down.sql", version, name)
	if err := ioutil.WriteFile(filepath.Join(path, down), nil, 0644); err != nil {
		return errors.Annotate(err, "writing down migration file failed")
	}

	fmt.Println("Version", version, "migration files created in", path)
	fmt.Println(up)
	fmt.Println(down)

	return nil
}
