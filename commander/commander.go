package commander

import (
	"strconv"
	"time"

	"github.com/juju/errors"
	"github.com/urfave/cli"
	"github.com/wallester/migrate/direction"
	"github.com/wallester/migrate/flag"
	"github.com/wallester/migrate/migrator"
)

// ICommander represents app commands
type ICommander interface {
	Create(c *cli.Context) error
	Up(c *cli.Context) error
	Down(c *cli.Context) error
}

type Commander struct {
	m migrator.IMigrator
}

var _ ICommander = (*Commander)(nil)

// New returns new instance
func New(m migrator.IMigrator) *Commander {
	return &Commander{
		m: m,
	}
}

// Create creates new migration files
func (cmd *Commander) Create(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		return errors.New("please specify migration name")
	}

	path := flag.Get(c, flag.Path)
	if path == "" {
		return flag.NewRequiredFlagError(flag.Path)
	}

	verbose := flag.GetBool(c, flag.Verbose)

	if _, err := cmd.m.Create(name, path, verbose); err != nil {
		return errors.Annotate(err, "creating migration failed")
	}

	return nil
}

// Up migrates up
func (cmd *Commander) Up(c *cli.Context) error {
	args, err := parseMigrateArguments(c)
	if err != nil {
		return errors.Annotate(err, "parsing parameters failed")
	}

	args.Direction = direction.Up
	if err := cmd.m.Migrate(*args); err != nil {
		return errors.Annotate(err, "migrating up failed")
	}

	return nil
}

// Down migrates down
func (cmd *Commander) Down(c *cli.Context) error {
	args, err := parseMigrateArguments(c)
	if err != nil {
		return errors.Annotate(err, "parsing parameters failed")
	}

	if args.Steps < 1 {
		return flag.NewRequiredFlagError("<n>")
	}

	args.Direction = direction.Down
	if err := cmd.m.Migrate(*args); err != nil {
		return errors.Annotate(err, "migrating down failed")
	}

	return nil
}

// private

func parseMigrateArguments(c *cli.Context) (*migrator.Args, error) {
	path := flag.Get(c, flag.Path)
	if path == "" {
		return nil, flag.NewRequiredFlagError(flag.Path)
	}

	url := flag.Get(c, flag.URL)
	if url == "" {
		return nil, flag.NewRequiredFlagError(flag.URL)
	}

	timeoutDuration := time.Second
	if s := flag.Get(c, flag.Timeout); s != "" {
		timeoutSeconds, err := strconv.Atoi(s)
		if err != nil {
			return nil, flag.NewWrongFormatFlagError(flag.Timeout)
		}

		timeoutDuration = time.Duration(timeoutSeconds) * time.Second
	}

	if s := flag.Get(c, flag.TimeoutDuration); s != "" {
		var err error
		timeoutDuration, err = time.ParseDuration(s)
		if err != nil {
			return nil, flag.NewWrongFormatFlagError(flag.TimeoutDuration)
		}
	}

	dbConnectionTimeoutDuration := time.Second
	if s := flag.Get(c, flag.DBConnectionTimeoutDuration); s != "" {
		var err error
		dbConnectionTimeoutDuration, err = time.ParseDuration(s)
		if err != nil {
			return nil, flag.NewWrongFormatFlagError(flag.DBConnectionTimeoutDuration)
		}
	}

	var steps int
	if s := c.Args().First(); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, flag.NewWrongFormatFlagError("<n>")
		}

		steps = n
	}

	noVerify := flag.GetBool(c, flag.NoVerify)
	verbose := flag.GetBool(c, flag.Verbose)

	return &migrator.Args{
		Path:                        path,
		URL:                         url,
		Steps:                       steps,
		NoVerify:                    noVerify,
		TimeoutDuration:             timeoutDuration,
		DBConnectionTimeoutDuration: dbConnectionTimeoutDuration,
		Verbose:                     verbose,
	}, nil
}
