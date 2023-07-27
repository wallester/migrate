package commander

import (
	"flag"
	"testing"
	"time"

	"github.com/juju/errors"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli"
	"github.com/wallester/migrate/direction"
	"github.com/wallester/migrate/file"
	"github.com/wallester/migrate/migrator"
)

type CommanderTestSuite struct {
	suite.Suite
	migratorMock *migrator.Mock
	commander    ICommander
	expectedErr  error
	flagSet      *flag.FlagSet
	ctx          *cli.Context
}

func (suite *CommanderTestSuite) SetupTest() {
	suite.migratorMock = &migrator.Mock{}
	suite.commander = New(suite.migratorMock)
	suite.expectedErr = errors.New("failure")
	suite.flagSet = flag.NewFlagSet("test", 0)
	suite.ctx = cli.NewContext(nil, suite.flagSet, nil)
}

func Test_Commander_TestSuite(t *testing.T) {
	suite.Run(t, &CommanderTestSuite{})
}

func (suite *CommanderTestSuite) Test_New_ReturnsInstance_InCaseOfSuccess() {
	// Act
	cmd := New(suite.migratorMock)

	// Assert
	suite.NotNil(cmd)
	suite.NotNil(cmd.Create)
	suite.NotNil(cmd.Up)
	suite.NotNil(cmd.Down)
}

func (suite *CommanderTestSuite) Test_Create_ReturnsError_InCaseOfMissingArgument() {
	// Act
	err := suite.commander.Create(suite.ctx)

	// Assert
	suite.EqualError(err, "please specify migration name")
}

func (suite *CommanderTestSuite) Test_Create_ReturnsError_InCaseOfMissingFlag() {
	// Arrange
	suite.Require().NoError(suite.flagSet.Parse([]string{"create_table_users"}))

	// Act
	err := suite.commander.Create(suite.ctx)

	// Assert
	suite.EqualError(err, "please specify path")
}

func (suite *CommanderTestSuite) Test_Create_ReturnsError_InCaseOfMigratorError() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.Require().NoError(suite.flagSet.Parse([]string{"--path", "testdata", "create_table_users"}))

	pair := &file.Pair{}
	suite.migratorMock.On("Create", "create_table_users", "testdata", false).Return(pair, suite.expectedErr).Once()

	// Act
	err := suite.commander.Create(suite.ctx)

	// Assert
	suite.EqualError(err, "creating migration failed: failure")
}

func (suite *CommanderTestSuite) Test_Create_ReturnsNil_InCaseOfSuccess() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.Require().NoError(suite.flagSet.Parse([]string{"--path", "testdata", "create_table_users"}))

	pair := &file.Pair{}
	suite.migratorMock.On("Create", "create_table_users", "testdata", false).Return(pair, nil).Once()

	// Act
	err := suite.commander.Create(suite.ctx)

	// Assert
	suite.NoError(err)
}

func (suite *CommanderTestSuite) Test_Up_ReturnsError_InCaseOfMissingPath() {
	// Act
	err := suite.commander.Up(suite.ctx)

	// Assert
	suite.EqualError(errors.Cause(err), "please specify path")
}

func (suite *CommanderTestSuite) Test_Up_ReturnsError_InCaseOfMissingURL() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.Require().NoError(suite.flagSet.Parse([]string{"--path", "testdata"}))

	// Act
	err := suite.commander.Up(suite.ctx)

	// Assert
	suite.EqualError(errors.Cause(err), "please specify url")
}

func (suite *CommanderTestSuite) Test_Up_ReturnsError_InCaseOfMigratorError() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	suite.flagSet.String("timeout-duration", "", "")
	suite.flagSet.Duration("db-conn-timeout-duration", 10*time.Second, "")
	suite.Require().NoError(
		suite.flagSet.Parse([]string{
			"--path", "testdata",
			"--url", "connectionurl",
			"--db-conn-timeout-duration", "10s",
			"--timeout-duration", "10s",
		}),
	)

	args := migrator.Args{
		Path:                        "testdata",
		URL:                         "connectionurl",
		Direction:                   direction.Up,
		TimeoutDuration:             10 * time.Second,
		DBConnectionTimeoutDuration: 10 * time.Second,
	}

	suite.migratorMock.On("Migrate", args).Return(suite.expectedErr).Once()

	// Act
	err := suite.commander.Up(suite.ctx)

	// Assert
	suite.EqualError(err, "migrating up failed: failure")
}

func (suite *CommanderTestSuite) Test_Up_ReturnsError_InCaseOfInvalidArgument() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	suite.Require().NoError(suite.flagSet.Parse([]string{"--path", "testdata", "--url", "connectionurl", "foobar"}))

	// Act
	err := suite.commander.Up(suite.ctx)

	// Assert
	suite.EqualError(errors.Cause(err), "parsing <n> failed")
}

func (suite *CommanderTestSuite) Test_Up_ReturnsNil_InCaseOfSuccessAndTimeoutDuration() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	suite.flagSet.String("timeout-duration", "", "")
	suite.flagSet.Duration("db-conn-timeout-duration", 10*time.Second, "")
	suite.Require().NoError(
		suite.flagSet.Parse([]string{
			"--path", "testdata",
			"--url", "connectionurl",
			"--db-conn-timeout-duration", "10s",
			"--timeout-duration", "10s",
		}),
	)

	args := migrator.Args{
		Path:                        "testdata",
		URL:                         "connectionurl",
		Direction:                   direction.Up,
		Steps:                       0,
		TimeoutDuration:             10 * time.Second,
		DBConnectionTimeoutDuration: 10 * time.Second,
	}

	suite.migratorMock.On("Migrate", args).Return(nil).Once()

	// Act
	err := suite.commander.Up(suite.ctx)

	// Assert
	suite.NoError(err)
}

func (suite *CommanderTestSuite) Test_Up_ReturnsNil_InCaseOfSuccessAndTimeout() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	suite.flagSet.String("timeout", "", "")
	suite.flagSet.Duration("db-conn-timeout-duration", 10*time.Second, "")
	suite.Require().NoError(
		suite.flagSet.Parse([]string{
			"--path", "testdata",
			"--url", "connectionurl",
			"--db-conn-timeout-duration", "10s",
			"--timeout", "10",
		}),
	)

	args := migrator.Args{
		Path:                        "testdata",
		URL:                         "connectionurl",
		Direction:                   direction.Up,
		Steps:                       0,
		TimeoutDuration:             10 * time.Second,
		DBConnectionTimeoutDuration: 10 * time.Second,
	}

	suite.migratorMock.On("Migrate", args).Return(nil).Once()

	// Act
	err := suite.commander.Up(suite.ctx)

	// Assert
	suite.NoError(err)
}

func (suite *CommanderTestSuite) Test_Up_ReturnsNil_InCaseOfArgumentN() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	suite.flagSet.String("timeout-duration", "", "")
	suite.flagSet.Duration("db-conn-timeout-duration", 10*time.Second, "")
	suite.Require().NoError(
		suite.flagSet.Parse([]string{
			"--path", "testdata",
			"--url", "connectionurl",
			"--db-conn-timeout-duration", "10s",
			"--timeout-duration", "10s",
			"10",
		}),
	)

	args := migrator.Args{
		Path:                        "testdata",
		URL:                         "connectionurl",
		Direction:                   direction.Up,
		Steps:                       10,
		TimeoutDuration:             10 * time.Second,
		DBConnectionTimeoutDuration: 10 * time.Second,
	}

	suite.migratorMock.On("Migrate", args).Return(nil).Once()

	// Act
	err := suite.commander.Up(suite.ctx)

	// Assert
	suite.NoError(err)
}

func (suite *CommanderTestSuite) Test_Down_ReturnsError_InCaseOfMissingPath() {
	// Act
	err := suite.commander.Down(suite.ctx)

	// Assert
	suite.EqualError(errors.Cause(err), "please specify path")
}

func (suite *CommanderTestSuite) Test_Down_ReturnsError_InCaseOfMissingURL() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.Require().NoError(suite.flagSet.Parse([]string{"--path", "testdata"}))

	// Act
	err := suite.commander.Down(suite.ctx)

	// Assert
	suite.EqualError(errors.Cause(err), "please specify url")
}

func (suite *CommanderTestSuite) Test_Down_ReturnsError_InCaseOfMigratorError() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	suite.flagSet.String("timeout-duration", "", "")
	suite.flagSet.Duration("db-conn-timeout-duration", 10*time.Second, "")
	suite.Require().NoError(
		suite.flagSet.Parse([]string{
			"--path", "testdata",
			"--url", "connectionurl",
			"--db-conn-timeout-duration", "10s",
			"--timeout-duration", "10s",
			"123",
		}),
	)

	args := migrator.Args{
		Path:                        "testdata",
		URL:                         "connectionurl",
		Direction:                   direction.Down,
		Steps:                       123,
		TimeoutDuration:             10 * time.Second,
		DBConnectionTimeoutDuration: 10 * time.Second,
	}

	suite.migratorMock.On("Migrate", args).Return(suite.expectedErr).Once()

	// Act
	err := suite.commander.Down(suite.ctx)

	// Assert
	suite.EqualError(err, "migrating down failed: failure")
}

func (suite *CommanderTestSuite) Test_Down_ReturnsError_InCaseOfMissingArgumentN() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	suite.Require().NoError(suite.flagSet.Parse([]string{"--path", "testdata", "--url", "connectionurl"}))

	// Act
	err := suite.commander.Down(suite.ctx)

	// Assert
	suite.EqualError(err, "please specify <n>")
}

func (suite *CommanderTestSuite) Test_Down_ReturnsNil_InCaseOfSuccessAndTimeoutDuration() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	suite.flagSet.String("timeout-duration", "", "")
	suite.flagSet.Duration("db-conn-timeout-duration", 10*time.Second, "")
	suite.Require().NoError(
		suite.flagSet.Parse([]string{
			"--path", "testdata",
			"--url", "connectionurl",
			"--db-conn-timeout-duration", "10s",
			"--timeout-duration", "10s",
			"123",
		}),
	)

	args := migrator.Args{
		Path:                        "testdata",
		URL:                         "connectionurl",
		Direction:                   direction.Down,
		Steps:                       123,
		TimeoutDuration:             10 * time.Second,
		DBConnectionTimeoutDuration: 10 * time.Second,
	}

	suite.migratorMock.On("Migrate", args).Return(nil).Once()

	// Act
	err := suite.commander.Down(suite.ctx)

	// Assert
	suite.NoError(err)
}

func (suite *CommanderTestSuite) Test_Down_ReturnsNil_InCaseOfSuccessAndTimeout() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	suite.flagSet.String("timeout", "", "")
	suite.flagSet.Duration("db-conn-timeout-duration", 10*time.Second, "")
	suite.Require().NoError(
		suite.flagSet.Parse([]string{
			"--path", "testdata",
			"--url", "connectionurl",
			"--db-conn-timeout-duration", "10s",
			"--timeout", "10",
			"123",
		}),
	)

	args := migrator.Args{
		Path:                        "testdata",
		URL:                         "connectionurl",
		Direction:                   direction.Down,
		Steps:                       123,
		TimeoutDuration:             10 * time.Second,
		DBConnectionTimeoutDuration: 10 * time.Second,
	}

	suite.migratorMock.On("Migrate", args).Return(nil).Once()

	// Act
	err := suite.commander.Down(suite.ctx)

	// Assert
	suite.NoError(err)
}
