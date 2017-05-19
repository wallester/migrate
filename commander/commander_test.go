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
	flagSet      *flag.FlagSet
	ctx          *cli.Context
}

func (suite *CommanderTestSuite) SetupTest() {
	suite.migratorMock = &migrator.Mock{}
	suite.instance = New(suite.migratorMock)
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
	err := suite.instance.Create(suite.ctx)

	// Assert
	suite.NotNil(err)
	suite.EqualError(err, "please specify migration name")
}

func (suite *CommanderTestSuite) Test_Create_ReturnsError_InCaseOfMissingFlag() {
	// Arrange
	if err := suite.flagSet.Parse([]string{"create_table_users"}); err != nil {
		suite.FailNow(err.Error())
	}

	// Act
	err := suite.instance.Create(suite.ctx)

	// Assert
	suite.NotNil(err)
	suite.EqualError(err, "please specify path")
}

func (suite *CommanderTestSuite) Test_Create_ReturnsError_InCaseOfMigratorError() {
	// Arrange
	suite.flagSet.String("path", "", "")
	if err := suite.flagSet.Parse([]string{"--path", "testdata", "create_table_users"}); err != nil {
		suite.FailNow(err.Error())
	}
	pair := &file.Pair{}
	suite.migratorMock.On("Create", "create_table_users", "testdata").Return(pair, suite.expectedErr).Once()

	// Act
	err := suite.instance.Create(suite.ctx)

	// Assert
	suite.NotNil(err)
	suite.EqualError(err, "creating migration failed: failure")
}

func (suite *CommanderTestSuite) Test_Create_ReturnsNil_InCaseOfSuccess() {
	// Arrange
	suite.flagSet.String("path", "", "")
	if err := suite.flagSet.Parse([]string{"--path", "testdata", "create_table_users"}); err != nil {
		suite.FailNow(err.Error())
	}
	pair := &file.Pair{}
	suite.migratorMock.On("Create", "create_table_users", "testdata").Return(pair, nil).Once()

	// Act
	err := suite.instance.Create(suite.ctx)

	// Assert
	suite.Nil(err)
}

func (suite *CommanderTestSuite) Test_Up_ReturnError_InCaseOfMissingPath() {
	// Act
	err := suite.instance.Up(suite.ctx)

	// Assert
	suite.NotNil(err)
	suite.EqualError(errors.Cause(err), "please specify path")
}

func (suite *CommanderTestSuite) Test_Up_ReturnError_InCaseOfMissingURL() {
	// Arrange
	suite.flagSet.String("path", "", "")
	if err := suite.flagSet.Parse([]string{"--path", "testdata"}); err != nil {
		suite.FailNow(err.Error())
	}

	// Act
	err := suite.instance.Up(suite.ctx)

	// Assert
	suite.NotNil(err)
	suite.EqualError(errors.Cause(err), "please specify url")
}

func (suite *CommanderTestSuite) Test_Up_ReturnError_InCaseOfMigratorError() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	if err := suite.flagSet.Parse([]string{"--path", "testdata", "--url", "connectionurl"}); err != nil {
		suite.FailNow(err.Error())
	}
	suite.migratorMock.On("Migrate", "testdata", "connectionurl", true, 0).Return(suite.expectedErr).Once()

	// Act
	err := suite.instance.Up(suite.ctx)

	// Assert
	suite.NotNil(err)
	suite.EqualError(err, "migrating up failed: failure")
}

func (suite *CommanderTestSuite) Test_Up_ReturnError_InCaseOfInvalidArgument() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	if err := suite.flagSet.Parse([]string{"--path", "testdata", "--url", "connectionurl", "foobar"}); err != nil {
		suite.FailNow(err.Error())
	}

	// Act
	err := suite.instance.Up(suite.ctx)

	// Assert
	suite.NotNil(err)
	suite.EqualError(errors.Cause(err), "strconv.Atoi: parsing \"foobar\": invalid syntax")
}

func (suite *CommanderTestSuite) Test_Up_ReturnNil_InCaseOfSuccess() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	if err := suite.flagSet.Parse([]string{"--path", "testdata", "--url", "connectionurl"}); err != nil {
		suite.FailNow(err.Error())
	}
	suite.migratorMock.On("Migrate", "testdata", "connectionurl", true, 0).Return(nil).Once()

	// Act
	err := suite.instance.Up(suite.ctx)

	// Assert
	suite.Nil(err)
}

func (suite *CommanderTestSuite) Test_Up_ReturnNil_InCaseOfArgumentN() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	if err := suite.flagSet.Parse([]string{"--path", "testdata", "--url", "connectionurl", "10"}); err != nil {
		suite.FailNow(err.Error())
	}
	suite.migratorMock.On("Migrate", "testdata", "connectionurl", true, 10).Return(nil).Once()

	// Act
	err := suite.instance.Up(suite.ctx)

	// Assert
	suite.Nil(err)
}

func (suite *CommanderTestSuite) Test_Down_ReturnError_InCaseOfMissingPath() {
	// Act
	err := suite.instance.Down(suite.ctx)

	// Assert
	suite.NotNil(err)
	suite.EqualError(errors.Cause(err), "please specify path")
}

func (suite *CommanderTestSuite) Test_Down_ReturnError_InCaseOfMissingURL() {
	// Arrange
	suite.flagSet.String("path", "", "")
	if err := suite.flagSet.Parse([]string{"--path", "testdata"}); err != nil {
		suite.FailNow(err.Error())
	}

	// Act
	err := suite.instance.Down(suite.ctx)

	// Assert
	suite.NotNil(err)
	suite.EqualError(errors.Cause(err), "please specify url")
}

func (suite *CommanderTestSuite) Test_Down_ReturnError_InCaseOfMigratorError() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	if err := suite.flagSet.Parse([]string{"--path", "testdata", "--url", "connectionurl"}); err != nil {
		suite.FailNow(err.Error())
	}
	suite.migratorMock.On("Migrate", "testdata", "connectionurl", false, 0).Return(suite.expectedErr).Once()

	// Act
	err := suite.instance.Down(suite.ctx)

	// Assert
	suite.NotNil(err)
	suite.EqualError(err, "migrating down failed: failure")
}

func (suite *CommanderTestSuite) Test_Down_ReturnNil_InCaseOfSuccess() {
	// Arrange
	suite.flagSet.String("path", "", "")
	suite.flagSet.String("url", "", "")
	if err := suite.flagSet.Parse([]string{"--path", "testdata", "--url", "connectionurl"}); err != nil {
		suite.FailNow(err.Error())
	}
	suite.migratorMock.On("Migrate", "testdata", "connectionurl", false, 0).Return(nil).Once()

	// Act
	err := suite.instance.Down(suite.ctx)

	// Assert
	suite.Nil(err)
}
