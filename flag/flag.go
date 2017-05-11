package flag

import (
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
		Usage:  "migrations folder, defaults to current working directory",
		EnvVar: "MIGRATE_PATH",
	}
	// Path represents a cli flag
	Path = cli.StringFlag{
		Name:   FlagPath,
		Usage:  "database URL, for example postgres://user@host:port/database",
		EnvVar: "MIGRATE_URL",
	}
)
