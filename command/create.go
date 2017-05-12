package command

import (
	"github.com/juju/errors"
	"github.com/urfave/cli"
	"github.com/wallester/migrate/flag"
	"github.com/wallester/migrate/migration"
)

// Create creates new migration files
func Create(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		return errors.New("please specify migration name")
	}

	path := flag.Get(c, flag.FlagPath)
	if path == "" {
		return flag.NewRequiredFlagError(flag.FlagPath)
	}

	if err := migration.Create(name, path); err != nil {
		return errors.Annotate(err, "creating migration failed")
	}

	return nil
}
