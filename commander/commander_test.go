package commander

import (
	"flag"
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"github.com/wallester/migrate/driver"
	"github.com/wallester/migrate/migrator"
	"github.com/wallester/migrate/printer"
)

func Test_New_ReturnsInstance_InCaseOfSuccess(t *testing.T) {
	// Act
	cmd := New(migrator.New(driver.New(), printer.New()))

	// Assert
	assert.NotNil(t, cmd)
	assert.NotNil(t, cmd.Create)
	assert.NotNil(t, cmd.Up)
	assert.NotNil(t, cmd.Down)
}

func Test_Create_ReturnsError_InCaseOfMissingArgument(t *testing.T) {
	// Arrange
	cmd := New(migrator.New(driver.New(), printer.New()))
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	// Act
	err := cmd.Create(c)

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, err, "please specify migration name")
}

func Test_Create_ReturnsError_InCaseOfMissingFlag(t *testing.T) {
	// Arrange
	cmd := New(migrator.New(driver.New(), printer.New()))
	set := flag.NewFlagSet("test", 0)
	if err := set.Parse([]string{"create_table_users"}); err != nil {
		assert.FailNow(t, err.Error())
	}
	c := cli.NewContext(nil, set, nil)

	// Act
	err := cmd.Create(c)

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, err, "please specify path")
}

func Test_Create_ReturnsError_InCaseOfMigratorError(t *testing.T) {
	// Arrange
	migratorMock := &migrator.Mock{}
	cmd := New(migratorMock)
	set := flag.NewFlagSet("test", 0)
	set.String("path", "", "")
	if err := set.Parse([]string{"--path", "testdata", "create_table_users"}); err != nil {
		assert.FailNow(t, err.Error())
	}
	c := cli.NewContext(nil, set, nil)
	expectedErr := errors.New("failure")
	migratorMock.On("Create", "create_table_users", "testdata").Return(expectedErr).Once()

	// Act
	err := cmd.Create(c)

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, err, "creating migration failed: failure")
}

func Test_Create_ReturnsNil_InCaseOfSuccess(t *testing.T) {
	// Arrange
	migratorMock := &migrator.Mock{}
	cmd := New(migratorMock)
	set := flag.NewFlagSet("test", 0)
	set.String("path", "", "")
	if err := set.Parse([]string{"--path", "testdata", "create_table_users"}); err != nil {
		assert.FailNow(t, err.Error())
	}
	c := cli.NewContext(nil, set, nil)
	migratorMock.On("Create", "create_table_users", "testdata").Return(nil).Once()

	// Act
	err := cmd.Create(c)

	// Assert
	assert.Nil(t, err)
}

func Test_Up_ReturnError_InCaseOfMissingPath(t *testing.T) {
	// Arrange
	cmd := New(migrator.New(driver.New(), printer.New()))
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	// Act
	err := cmd.Up(c)

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, err, "please specify path")
}

func Test_Up_ReturnError_InCaseOfMissingURL(t *testing.T) {
	// Arrange
	cmd := New(migrator.New(driver.New(), printer.New()))
	set := flag.NewFlagSet("test", 0)
	set.String("path", "", "")
	if err := set.Parse([]string{"--path", "testdata"}); err != nil {
		assert.FailNow(t, err.Error())
	}
	c := cli.NewContext(nil, set, nil)

	// Act
	err := cmd.Up(c)

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, err, "please specify url")
}

func Test_Up_ReturnError_InCaseOfMigratorError(t *testing.T) {
	// Arrange
	migratorMock := &migrator.Mock{}
	cmd := New(migratorMock)
	set := flag.NewFlagSet("test", 0)
	set.String("path", "", "")
	set.String("url", "", "")
	if err := set.Parse([]string{"--path", "testdata", "--url", "connectionurl"}); err != nil {
		assert.FailNow(t, err.Error())
	}
	c := cli.NewContext(nil, set, nil)
	expectedErr := errors.New("failure")
	migratorMock.On("Up", "testdata", "connectionurl").Return(expectedErr).Once()

	// Act
	err := cmd.Up(c)

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, err, "migrating up failed: failure")
}

func Test_Up_ReturnNil_InCaseOfSuccess(t *testing.T) {
	// Arrange
	migratorMock := &migrator.Mock{}
	cmd := New(migratorMock)
	set := flag.NewFlagSet("test", 0)
	set.String("path", "", "")
	set.String("url", "", "")
	if err := set.Parse([]string{"--path", "testdata", "--url", "connectionurl"}); err != nil {
		assert.FailNow(t, err.Error())
	}
	c := cli.NewContext(nil, set, nil)
	migratorMock.On("Up", "testdata", "connectionurl").Return(nil).Once()

	// Act
	err := cmd.Up(c)

	// Assert
	assert.Nil(t, err)
}

func Test_Down_ReturnError_InCaseOfMissingPath(t *testing.T) {
	// Arrange
	cmd := New(migrator.New(driver.New(), printer.New()))
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	// Act
	err := cmd.Down(c)

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, err, "please specify path")
}

func Test_Down_ReturnError_InCaseOfMissingURL(t *testing.T) {
	// Arrange
	cmd := New(migrator.New(driver.New(), printer.New()))
	set := flag.NewFlagSet("test", 0)
	set.String("path", "", "")
	if err := set.Parse([]string{"--path", "testdata"}); err != nil {
		assert.FailNow(t, err.Error())
	}
	c := cli.NewContext(nil, set, nil)

	// Act
	err := cmd.Down(c)

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, err, "please specify url")
}

func Test_Down_ReturnError_InCaseOfMigratorError(t *testing.T) {
	// Arrange
	migratorMock := &migrator.Mock{}
	cmd := New(migratorMock)
	set := flag.NewFlagSet("test", 0)
	set.String("path", "", "")
	set.String("url", "", "")
	if err := set.Parse([]string{"--path", "testdata", "--url", "connectionurl"}); err != nil {
		assert.FailNow(t, err.Error())
	}
	c := cli.NewContext(nil, set, nil)
	expectedErr := errors.New("failure")
	migratorMock.On("Down", "testdata", "connectionurl").Return(expectedErr).Once()

	// Act
	err := cmd.Down(c)

	// Assert
	assert.NotNil(t, err)
	assert.EqualError(t, err, "migrating down failed: failure")
}

func Test_Down_ReturnNil_InCaseOfSuccess(t *testing.T) {
	// Arrange
	migratorMock := &migrator.Mock{}
	cmd := New(migratorMock)
	set := flag.NewFlagSet("test", 0)
	set.String("path", "", "")
	set.String("url", "", "")
	if err := set.Parse([]string{"--path", "testdata", "--url", "connectionurl"}); err != nil {
		assert.FailNow(t, err.Error())
	}
	c := cli.NewContext(nil, set, nil)
	migratorMock.On("Down", "testdata", "connectionurl").Return(nil).Once()

	// Act
	err := cmd.Down(c)

	// Assert
	assert.Nil(t, err)
}
