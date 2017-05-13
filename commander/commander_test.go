package commander

import (
	"flag"
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli"
	"github.com/wallester/migrate/file"
	"github.com/wallester/migrate/migrator"
)

type CommanderTestSuite struct {
	suite.Suite
	migratorMock *migrator.Mock
	instance     Commander
	expectedErr  error
	set          *flag.FlagSet
	c            *cli.Context
}

func (suite *CommanderTestSuite) SetupTest() {
	suite.migratorMock = &migrator.Mock{}
	suite.instance = New(suite.migratorMock)
	suite.expectedErr = errors.New("failure")
	suite.set = flag.NewFlagSet("test", 0)
	suite.c = cli.NewContext(nil, suite.set, nil)
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
	err := suite.instance.Create(suite.c)

	// Assert
	suite.NotNil(err)
	suite.EqualError(err, "please specify migration name")
}

func (suite *CommanderTestSuite) Test_Create_ReturnsError_InCaseOfMissingFlag() {
	// Arrange
	if err := suite.set.Parse([]string{"create_table_users"}); err != nil {
		suite.FailNow(err.Error())
	}

	// Act
	err := suite.instance.Create(suite.c)

	// Assert
	suite.NotNil(err)
	suite.EqualError(err, "please specify path")
}

func (suite *CommanderTestSuite) Test_Create_ReturnsError_InCaseOfMigratorError() {
	// Arrange
	suite.set.String("path", "", "")
	if err := suite.set.Parse([]string{"--path", "testdata", "create_table_users"}); err != nil {
		suite.FailNow(err.Error())
	}
	pair := &file.Pair{}
	suite.migratorMock.On("Create", "create_table_users", "testdata").Return(pair, suite.expectedErr).Once()

	// Act
	err := suite.instance.Create(suite.c)

	// Assert
	suite.NotNil(err)
	suite.EqualError(err, "creating migration failed: failure")
}

func (suite *CommanderTestSuite) Test_Create_ReturnsNil_InCaseOfSuccess() {
	// Arrange
	suite.set.String("path", "", "")
	if err := suite.set.Parse([]string{"--path", "testdata", "create_table_users"}); err != nil {
		suite.FailNow(err.Error())
	}
	pair := &file.Pair{}
	suite.migratorMock.On("Create", "create_table_users", "testdata").Return(pair, nil).Once()

	// Act
	err := suite.instance.Create(suite.c)

	// Assert
	suite.Nil(err)
}

func (suite *CommanderTestSuite) Test_Up_ReturnError_InCaseOfMissingPath() {
	// Act
	err := suite.instance.Up(suite.c)

	// Assert
	suite.NotNil(err)
	suite.EqualError(errors.Cause(err), "please specify path")
}

func (suite *CommanderTestSuite) Test_Up_ReturnError_InCaseOfMissingURL() {
	// Arrange
	suite.set.String("path", "", "")
	if err := suite.set.Parse([]string{"--path", "testdata"}); err != nil {
		suite.FailNow(err.Error())
	}

	// Act
	err := suite.instance.Up(suite.c)

	// Assert
	suite.NotNil(err)
	suite.EqualError(errors.Cause(err), "please specify url")
}

func (suite *CommanderTestSuite) Test_Up_ReturnError_InCaseOfMigratorError() {
	// Arrange
	suite.set.String("path", "", "")
	suite.set.String("url", "", "")
	if err := suite.set.Parse([]string{"--path", "testdata", "--url", "connectionurl"}); err != nil {
		suite.FailNow(err.Error())
	}
	suite.migratorMock.On("Migrate", "testdata", "connectionurl", true, 0).Return(suite.expectedErr).Once()

	// Act
	err := suite.instance.Up(suite.c)

	// Assert
	suite.NotNil(err)
	suite.EqualError(err, "migrating up failed: failure")
}

func (suite *CommanderTestSuite) Test_Up_ReturnNil_InCaseOfSuccess() {
	// Arrange
	suite.set.String("path", "", "")
	suite.set.String("url", "", "")
	if err := suite.set.Parse([]string{"--path", "testdata", "--url", "connectionurl"}); err != nil {
		suite.FailNow(err.Error())
	}
	suite.migratorMock.On("Migrate", "testdata", "connectionurl", true, 0).Return(nil).Once()

	// Act
	err := suite.instance.Up(suite.c)

	// Assert
	suite.Nil(err)
}

func (suite *CommanderTestSuite) Test_Down_ReturnError_InCaseOfMissingPath() {
	// Act
	err := suite.instance.Down(suite.c)

	// Assert
	suite.NotNil(err)
	suite.EqualError(errors.Cause(err), "please specify path")
}

func (suite *CommanderTestSuite) Test_Down_ReturnError_InCaseOfMissingURL() {
	// Arrange
	suite.set.String("path", "", "")
	if err := suite.set.Parse([]string{"--path", "testdata"}); err != nil {
		suite.FailNow(err.Error())
	}

	// Act
	err := suite.instance.Down(suite.c)

	// Assert
	suite.NotNil(err)
	suite.EqualError(errors.Cause(err), "please specify url")
}

func (suite *CommanderTestSuite) Test_Down_ReturnError_InCaseOfMigratorError() {
	// Arrange
	suite.set.String("path", "", "")
	suite.set.String("url", "", "")
	if err := suite.set.Parse([]string{"--path", "testdata", "--url", "connectionurl"}); err != nil {
		suite.FailNow(err.Error())
	}
	suite.migratorMock.On("Migrate", "testdata", "connectionurl", false, 0).Return(suite.expectedErr).Once()

	// Act
	err := suite.instance.Down(suite.c)

	// Assert
	suite.NotNil(err)
	suite.EqualError(err, "migrating down failed: failure")
}

func (suite *CommanderTestSuite) Test_Down_ReturnNil_InCaseOfSuccess() {
	// Arrange
	suite.set.String("path", "", "")
	suite.set.String("url", "", "")
	if err := suite.set.Parse([]string{"--path", "testdata", "--url", "connectionurl"}); err != nil {
		suite.FailNow(err.Error())
	}
	suite.migratorMock.On("Migrate", "testdata", "connectionurl", false, 0).Return(nil).Once()

	// Act
	err := suite.instance.Down(suite.c)

	// Assert
	suite.Nil(err)
}
