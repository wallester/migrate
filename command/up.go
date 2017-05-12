package command

import (
	"github.com/juju/errors"
	"github.com/urfave/cli"
	"github.com/wallester/migrate/flag"
	"github.com/wallester/migrate/migration"
)

// Up migrates up
func Up(c *cli.Context) error {
	path := flag.Get(c, flag.FlagPath)
	if path == "" {
		return flag.NewRequiredFlagError(flag.FlagPath)
	}

	url := flag.Get(c, flag.FlagURL)
	if url == "" {
		return flag.NewRequiredFlagError(flag.FlagURL)
	}

	if err := migration.Up(path, url); err != nil {
		return errors.Annotate(err, "migrating up failed")
	}

	return nil
}
