package commander

import (
	"github.com/juju/errors"
	"github.com/urfave/cli"
	"github.com/wallester/migrate/flag"
	"github.com/wallester/migrate/migrator"
)

// Commander represents app commands
type Commander interface {
	Create(c *cli.Context) error
	Up(c *cli.Context) error
	Down(c *cli.Context) error
}

type commander struct {
	m migrator.Migrator
}

// New returns new instance
func New(m migrator.Migrator) Commander {
	return &commander{m}
}

// Create creates new migration files
func (cmd *commander) Create(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		return errors.New("please specify migration name")
	}

	path := flag.Get(c, flag.FlagPath)
	if path == "" {
		return flag.NewRequiredFlagError(flag.FlagPath)
	}

	if err := cmd.m.Create(name, path); err != nil {
		return errors.Annotate(err, "creating migration failed")
	}

	return nil
}

// Up migrates up
func (cmd *commander) Up(c *cli.Context) error {
	path := flag.Get(c, flag.FlagPath)
	if path == "" {
		return flag.NewRequiredFlagError(flag.FlagPath)
	}

	url := flag.Get(c, flag.FlagURL)
	if url == "" {
		return flag.NewRequiredFlagError(flag.FlagURL)
	}

	if err := cmd.m.Up(path, url); err != nil {
		return errors.Annotate(err, "migrating up failed")
	}

	return nil
}

// Down migrates down
func (cmd *commander) Down(c *cli.Context) error {
	path := flag.Get(c, flag.FlagPath)
	if path == "" {
		return flag.NewRequiredFlagError(flag.FlagPath)
	}

	url := flag.Get(c, flag.FlagURL)
	if url == "" {
		return flag.NewRequiredFlagError(flag.FlagURL)
	}

	if err := cmd.m.Down(path, url); err != nil {
		return errors.Annotate(err, "migrating down failed")
	}

	return nil
}
