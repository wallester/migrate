package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func Test_New_ReturnsInstance_InCaseOfSuccess(t *testing.T) {
	// Act
	app := New()

	// Assert
	assert.NotNil(t, app)
	assert.Equal(t, "migrate", app.Name)
	assert.Equal(t, "Command line tool for PostgreSQL migrations", app.Usage)
	assert.NotNil(t, getCommand("create", app.Commands))
	assert.NotNil(t, getCommand("up", app.Commands))
	assert.NotNil(t, getCommand("down", app.Commands))
}

func getCommand(name string, commands []cli.Command) *cli.Command {
	for _, command := range commands {
		if command.Name == name {
			return &command
		}
	}

	return nil
}
