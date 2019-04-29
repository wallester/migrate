package flag

import (
	"github.com/juju/errors"
	"github.com/urfave/cli"
)

const (
	// FlagURL represents database URL.
	FlagURL = "url"
	// FlagPath represents migrations path.
	FlagPath = "path"
	// FlagTimeout represents execution timeout in seconds.
	FlagTimeout = "timeout"
	// FlagNoVerify skips verification of already migrated older migrations.
	FlagNoVerify = "no-verify"
)

var (
	// URL represents a cli flag.
	URL = cli.StringFlag{
		Name:   FlagURL,
		Usage:  "database URL, for example postgres://user@host:port/database",
		EnvVar: "MIGRATE_URL",
	}

	// Path represents a cli flag.
	Path = cli.StringFlag{
		Name:   FlagPath,
		Usage:  "migrations folder, defaults to current working directory",
		EnvVar: "MIGRATE_PATH",
	}

	// Timeout represents a cli flag.
	Timeout = cli.StringFlag{
		Name:   FlagTimeout,
		Usage:  "execution timeout in seconds, defaults to 1 second",
		EnvVar: "MIGRATE_TIMEOUT",
	}

	// NoVerify represents a cli flag.
	NoVerify = cli.BoolFlag{
		Name:   FlagNoVerify,
		Usage:  "skip verification of already migrated older migrations",
		EnvVar: "MIGRATE_NO_VERIFY",
	}
)

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
