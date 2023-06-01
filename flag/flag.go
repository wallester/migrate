package flag

import (
	"github.com/juju/errors"
	"github.com/urfave/cli"
)

const (
	// URL represents database URL.
	URL = "url"
	// Path represents migrations path.
	Path = "path"
	// Timeout represents execution timeout in seconds.
	Timeout = "timeout"
	// NoVerify skips verification of already migrated older migrations.
	NoVerify = "no-verify"
)

var Flags = map[string]cli.Flag{
	URL: cli.StringFlag{
		Name:   URL,
		Usage:  "database URL, for example postgres://user@host:port/database",
		EnvVar: "MIGRATE_URL",
	},
	Path: cli.StringFlag{
		Name:   Path,
		Usage:  "migrations folder, defaults to current working directory",
		EnvVar: "MIGRATE_PATH",
	},
	Timeout: cli.StringFlag{
		Name:   Timeout,
		Usage:  "execution timeout in seconds, defaults to 1 second",
		EnvVar: "MIGRATE_TIMEOUT",
	},
	NoVerify: cli.BoolFlag{
		Name:   NoVerify,
		Usage:  "skip verification of already migrated older migrations",
		EnvVar: "MIGRATE_NO_VERIFY",
	},
}

// Get returns a flag value.
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

// GetBool returns a boolean flag value.
func GetBool(c *cli.Context, name string) bool {
	return c.Bool(name) || c.GlobalBool(name)
}

// NewRequiredFlagError returns new required flag error.
func NewRequiredFlagError(name string) error {
	return errors.New("please specify " + name)
}

// NewWrongFormatFlagError returns new wrong format flag error.
func NewWrongFormatFlagError(name string) error {
	return errors.New("parsing " + name + " failed")
}
