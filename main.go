package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/wallester/migrate/command"
)

const (
	flagURL  = "url"
	flagPath = "path"
)

func main() {
	app := cli.NewApp()
	app.Name = "migrate"
	app.Usage = "Command line tool for PostgreSQL migrations"
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "create",
			Usage:     "Create a new migration",
			ArgsUsage: "<name>",
			Action:    command.Create,
		},
		cli.Command{
			Name:   "up",
			Usage:  "Apply all -up- migrations",
			Action: command.Up,
		},
		cli.Command{
			Name:   "down",
			Usage:  "Apply all -down- migrations",
			Action: command.Down,
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   flagPath,
			Usage:  "defaults to current working directory",
			EnvVar: "MIGRATE_PATH",
		},
		cli.StringFlag{
			Name:   flagURL,
			Usage:  "postgres://user@host:port/database",
			EnvVar: "MIGRATE_URL",
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
