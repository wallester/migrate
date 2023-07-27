package app

import (
	"github.com/urfave/cli"
	"github.com/wallester/migrate/commander"
	"github.com/wallester/migrate/driver/postgres"
	"github.com/wallester/migrate/flag"
	"github.com/wallester/migrate/migrator"
	"github.com/wallester/migrate/printer"
)

// New returns new cli.App instance
func New() *cli.App {
	p := printer.New()
	d := postgres.New()
	m := migrator.New(d, p)
	cmd := commander.New(m)

	app := cli.NewApp()
	app.Name = "migrate"
	app.Usage = "Command line tool for Postgres migrations"
	app.Version = "1.0.2"
	app.Commands = []cli.Command{
		{
			Name:      "create",
			Usage:     "Create a new migration",
			ArgsUsage: "<name>",
			Action:    cmd.Create,
			Flags: []cli.Flag{
				flag.Flags[flag.Path],
				flag.Flags[flag.Verbose],
			},
		},
		{
			Name:      "up",
			Usage:     "Apply <n> or all up migrations",
			Action:    cmd.Up,
			ArgsUsage: "<n>",
			Flags: []cli.Flag{
				flag.Flags[flag.Path],
				flag.Flags[flag.URL],
				flag.Flags[flag.Timeout],
				flag.Flags[flag.TimeoutDuration],
				flag.Flags[flag.Verbose],
			},
		},
		{
			Name:      "down",
			Usage:     "Apply <n> down migration(s)",
			Action:    cmd.Down,
			ArgsUsage: "<n>",
			Flags: []cli.Flag{
				flag.Flags[flag.Path],
				flag.Flags[flag.URL],
				flag.Flags[flag.Timeout],
				flag.Flags[flag.TimeoutDuration],
				flag.Flags[flag.Verbose],
			},
		},
	}

	app.Flags = []cli.Flag{
		flag.Flags[flag.Path],
		flag.Flags[flag.URL],
		flag.Flags[flag.Timeout],
		flag.Flags[flag.TimeoutDuration],
		flag.Flags[flag.NoVerify],
		flag.Flags[flag.Verbose],
	}

	return app
}
