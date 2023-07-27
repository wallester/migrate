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
	// Timeout represents execution timeout in seconds. Default value: 1s.
	// Deprecated: use --timeout-duration instead.
	Timeout = "timeout"
	// TimeoutDuration represents execution timeout in duration.
	// TimeoutDuration will override timeout setting. Default value: 1s.
	TimeoutDuration = "timeout-duration"
	// DBConnectionTimeoutDuration represents database connection timeout in duration. Default value: 1s.
	DBConnectionTimeoutDuration = "db-conn-timeout-duration"
	// NoVerify skips verification of already migrated older migrations.
	NoVerify = "no-verify"
	// Verbose enables verbose output.
	Verbose = "verbose"
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
	// Deprecated, use TimeoutDuration instead.
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
	Verbose: cli.BoolFlag{
		Name:   Verbose,
		Usage:  "enable verbose output",
		EnvVar: "MIGRATE_VERBOSE",
	},
	TimeoutDuration: cli.DurationFlag{
		Name:   TimeoutDuration,
		Usage:  "execution timeout in duration, defaults to 1 second",
		EnvVar: "MIGRATE_TIMEOUT_DURATION",
	},
	DBConnectionTimeoutDuration: cli.DurationFlag{
		Name:   DBConnectionTimeoutDuration,
		Usage:  "database connection timeout in duration, defaults to 1 second",
		EnvVar: "MIGRATE_DB_CONN_TIMEOUT_DURATION",
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
