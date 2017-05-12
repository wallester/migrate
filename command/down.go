package command

import (
	"github.com/juju/errors"
	"github.com/urfave/cli"
	"github.com/wallester/migrate/flag"
	"github.com/wallester/migrate/migration"
)

// Down migrates down
func Down(c *cli.Context) error {
	path := flag.Get(c, flag.FlagPath)
	if path == "" {
		return flag.NewRequiredFlagError(flag.FlagPath)
	}

	url := flag.Get(c, flag.FlagURL)
	if url == "" {
		return flag.NewRequiredFlagError(flag.FlagURL)
	}

	if err := migration.Down(path, url); err != nil {
		return errors.Annotate(err, "migrating down failed")
	}

	return nil
}
