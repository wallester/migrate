package command

import (
	"github.com/juju/errors"
	"github.com/urfave/cli"
)

// Up migrates up
func Up(c *cli.Context) error {
	if err := migrate(c, true); err != nil {
		return errors.Annotate(err, "migrating up failed")
	}

	return nil
}
