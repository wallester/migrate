package command

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/urfave/cli"
	"github.com/wallester/migrate/flag"
)

// Create creates new migration files
func Create(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		return errors.New("please specify migration name")
	}

	name = strings.Replace(name, " ", "_", -1)

	path := ""
	if c.IsSet(flag.FlagPath) {
		path = c.String(flag.FlagPath)
	}

	if c.GlobalIsSet(flag.FlagPath) {
		path = c.GlobalString(flag.FlagPath)
	}

	if path == "" {
		return errors.New("please specify " + flag.Path.Name)
	}

	version := time.Now().Unix()

	up := fmt.Sprintf("%d_%s.up.sql", version, name)
	if err := ioutil.WriteFile(filepath.Join(path, up), nil, 0644); err != nil {
		return errors.Annotate(err, "writing up migration file failed")
	}

	down := fmt.Sprintf("%d_%s.down.sql", version, name)
	if err := ioutil.WriteFile(filepath.Join(path, down), nil, 0644); err != nil {
		return errors.Annotate(err, "writing down migration file failed")
	}

	fmt.Println("Version", version, "migration files created in", path)
	fmt.Println(up)
	fmt.Println(down)

	return nil
}
