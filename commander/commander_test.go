package commander

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wallester/migrate/driver"
	"github.com/wallester/migrate/migrator"
)

func Test_New_ReturnsInstance_InCaseOfSuccess(t *testing.T) {
	// Act
	cmd := New(migrator.New(driver.New()))

	// Assert
	assert.NotNil(t, cmd)
	assert.NotNil(t, cmd.Create)
	assert.NotNil(t, cmd.Up)
	assert.NotNil(t, cmd.Down)
}
