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
	cmd := commander.New(migrator.New(postgres.New(), printer.New()))

	app := cli.NewApp()
	app.Name = "migrate"
	app.Usage = "Command line tool for PostgreSQL migrations"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:      "create",
			Usage:     "Create a new migration",
			ArgsUsage: "<name>",
			Action:    cmd.Create,
			Flags: []cli.Flag{
				flag.Path,
			},
		},
		{
			Name:      "up",
			Usage:     "Apply <n> or all up migrations",
			Action:    cmd.Up,
			ArgsUsage: "<n>",
			Flags: []cli.Flag{
				flag.Path,
				flag.URL,
				flag.Timeout,
			},
		},
		{
			Name:      "down",
			Usage:     "Apply <n> or all down migrations",
			Action:    cmd.Down,
			ArgsUsage: "<n>",
			Flags: []cli.Flag{
				flag.Path,
				flag.URL,
				flag.Timeout,
			},
		},
	}
	app.Flags = []cli.Flag{
		flag.Path,
		flag.URL,
		flag.Timeout,
	}
	return app
}
