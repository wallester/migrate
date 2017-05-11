package command

import (
	"github.com/juju/errors"
	"github.com/urfave/cli"
)

// Down migrates down
func Down(c *cli.Context) error {
	if err := migrate(c, false); err != nil {
		return errors.Annotate(err, "migrating down failed")
	}

	return nil
}
