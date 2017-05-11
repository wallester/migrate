package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/wallester/migrate/command"
	"github.com/wallester/migrate/flag"
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
			Flags: []cli.Flag{
				flag.Path,
			},
		},
		cli.Command{
			Name:   "up",
			Usage:  "Apply all -up- migrations",
			Action: command.Up,
			Flags: []cli.Flag{
				flag.Path,
				flag.URL,
			},
		},
		cli.Command{
			Name:   "down",
			Usage:  "Apply all -down- migrations",
			Action: command.Down,
			Flags: []cli.Flag{
				flag.Path,
				flag.URL,
			},
		},
	}
	app.Flags = []cli.Flag{
		flag.Path,
		flag.URL,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
