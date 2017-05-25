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
	assert.NotNil(t, app.Commands)
	assert.True(t, hasCommand("create", app.Commands))
	assert.True(t, hasCommand("up", app.Commands))
	assert.True(t, hasCommand("down", app.Commands))
	assert.NotNil(t, app.Flags)
	assert.True(t, hasFlag("path", app.Flags))
	assert.True(t, hasFlag("url", app.Flags))
}

func hasCommand(name string, commands []cli.Command) bool {
	for _, command := range commands {
		if command.Name == name {
			return true
		}
	}

	return false
}

func hasFlag(name string, flags []cli.Flag) bool {
	for _, flag := range flags {
		if flag.GetName() == name {
			return true
		}
	}

	return false
}
