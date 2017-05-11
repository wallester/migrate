package flag

import (
	"github.com/juju/errors"
	"github.com/urfave/cli"
)

const (
	// FlagURL represents database URL
	FlagURL = "url"
	// FlagPath represents migrations path
	FlagPath = "path"
)

var (
	// URL represents a cli flag
	URL = cli.StringFlag{
		Name:   FlagURL,
		Usage:  "database URL, for example postgres://user@host:port/database",
		EnvVar: "MIGRATE_URL",
	}
	// Path represents a cli flag
	Path = cli.StringFlag{
		Name:   FlagPath,
		Usage:  "migrations folder, defaults to current working directory",
		EnvVar: "MIGRATE_PATH",
	}
)

// Get returns a flag value
func Get(c *cli.Context, name string) string {
	value := ""
	if c.IsSet(name) {
		value = c.String(name)
	}

	if c.GlobalIsSet(name) {
		value = c.GlobalString(name)
	}

	return value
}

// NewRequiredFlagError returns new required flag error
func NewRequiredFlagError(name string) error {
	return errors.New("please specify " + name)
}
