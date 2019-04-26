package commander

import (
	"strconv"

	"github.com/juju/errors"
	"github.com/urfave/cli"
	"github.com/wallester/migrate/flag"
	"github.com/wallester/migrate/migrator"
)

// Commander represents app commands
type Commander interface {
	Create(c *cli.Context) error
	Up(c *cli.Context) error
	Down(c *cli.Context) error
}

type commander struct {
	m migrator.Migrator
}

// New returns new instance
func New(m migrator.Migrator) Commander {
	return &commander{
		m: m,
	}
}

// Create creates new migration files
func (cmd *commander) Create(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		return errors.New("please specify migration name")
	}

	path := flag.Get(c, flag.FlagPath)
	if path == "" {
		return flag.NewRequiredFlagError(flag.FlagPath)
	}

	if _, err := cmd.m.Create(name, path); err != nil {
		return errors.Annotate(err, "creating migration failed")
	}

	return nil
}

// Up migrates up
func (cmd *commander) Up(c *cli.Context) error {
	args, err := parseMigrateArguments(c)
	if err != nil {
		return errors.Annotate(err, "parsing parameters failed")
	}

	args.Up = true
	if err := cmd.m.Migrate(*args); err != nil {
		return errors.Annotate(err, "migrating up failed")
	}

	return nil
}

// Down migrates down
func (cmd *commander) Down(c *cli.Context) error {
	args, err := parseMigrateArguments(c)
	if err != nil {
		return errors.Annotate(err, "parsing parameters failed")
	}

	if args.Steps < 1 {
		return flag.NewRequiredFlagError("<n>")
	}

	args.Up = false
	if err := cmd.m.Migrate(*args); err != nil {
		return errors.Annotate(err, "migrating down failed")
	}

	return nil
}

func parseMigrateArguments(c *cli.Context) (*migrator.Args, error) {
	path := flag.Get(c, flag.FlagPath)
	if path == "" {
		return nil, flag.NewRequiredFlagError(flag.FlagPath)
	}

	url := flag.Get(c, flag.FlagURL)
	if url == "" {
		return nil, flag.NewRequiredFlagError(flag.FlagURL)
	}

	var timeoutSeconds int
	s := flag.Get(c, flag.FlagTimeout)
	if s == "" {
		timeoutSeconds = 1
	} else {
		var err error
		timeoutSeconds, err = strconv.Atoi(s)
		if err != nil {
			return nil, flag.NewWrongFormatFlagError(flag.FlagTimeout)
		}
	}

	var steps int
	s = c.Args().First()
	if s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, flag.NewWrongFormatFlagError("<n>")
		}

		steps = n
	}

	return &migrator.Args{
		Path:           path,
		URL:            url,
		TimeoutSeconds: timeoutSeconds,
		Steps:          steps,
		NoVerify:       flag.GetBool(c, flag.FlagNoVerify),
	}, nil
}
