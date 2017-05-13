package migrator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/wallester/migrate/driver"
	"github.com/wallester/migrate/file"
	"github.com/wallester/migrate/printer"
)

type MigratorTestSuite struct {
	suite.Suite
	driverMock  *driver.Mock
	output      *printer.Recorder
	instance    Migrator
	expectedErr error
}

func (suite *MigratorTestSuite) SetupTest() {
	suite.driverMock = &driver.Mock{}
	suite.output = &printer.Recorder{}
	suite.instance = New(suite.driverMock, suite.output)
	suite.expectedErr = errors.New("failure")
}

func Test_Migrator_TestSuite(t *testing.T) {
	suite.Run(t, &MigratorTestSuite{})
}

func (suite *MigratorTestSuite) Test_New_ReturnsNewInstance_InCaseOfSuccess() {
	// Act
	instance := New(&driver.Mock{}, printer.New())

	// Assert
	assert.NotNil(suite.T(), instance)
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsNil_InCaseOfNoUpMigrationsToRun() {
	// Arrange
	// The following versions are from ../testdata.
	// We'll mark all of them as already migrated, meaning
	// no up migrations need to run.
	migrations := map[int64]bool{
		1494538273: true,
		1494538317: true,
		1494538407: true,
	}
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	suite.driverMock.On("SelectMigrations", mock.AnythingOfType("*context.timerCtx")).Return(migrations, nil).Once()
	suite.driverMock.On("Close").Once()

	// Act
	err := suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	assert.Nil(suite.T(), errors.Cause(err))
	assert.True(suite.T(), suite.output.Contains("seconds"))
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsError_InCaseOfDriverOpenError() {
	// Arrange
	suite.driverMock.On("Open", "connectionurl").Return(suite.expectedErr).Once()

	// Act
	err := suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "opening database connection failed: failure")
	assert.Empty(suite.T(), suite.output.String())
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsError_InCaseOfDriverCreateMigrationsTableError() {
	// Arrange
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(suite.expectedErr).Once()
	suite.driverMock.On("Close").Once()

	// Act
	err := suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "migrating failed: creating migrations table failed: failure")
	assert.Empty(suite.T(), suite.output.String())
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsErr_InCaseOfDriverSelectMigrationsError() {
	// Arrange
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	suite.driverMock.On("SelectMigrations", mock.AnythingOfType("*context.timerCtx")).Return(nil, suite.expectedErr).Once()
	suite.driverMock.On("Close").Once()

	// Act
	err := suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "migrating failed: selecting existing migrations failed: failure")
	assert.Empty(suite.T(), suite.output.String())
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsError_InCaseOfDriverApplyMigrationsError() {
	// Arrange
	// The following versions are from ../testdata.
	// We'll mark one of them as not migrated yet, meaning it needs
	// to be migrated up.
	migrations := map[int64]bool{
		1494538273: true,
		1494538317: false,
		1494538407: true,
	}
	files, err := file.ListFiles(filepath.Join("..", "testdata"), true)
	if err != nil {
		suite.FailNow(err.Error())
	}
	needsMigration := []file.File{
		*file.FindByVersion(1494538317, files),
	}
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	suite.driverMock.On("SelectMigrations", mock.AnythingOfType("*context.timerCtx")).Return(migrations, nil).Once()
	suite.driverMock.On("ApplyMigrations", mock.AnythingOfType("*context.timerCtx"), needsMigration, true).Return(suite.expectedErr).Once()
	suite.driverMock.On("Close").Once()

	// Act
	err = suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	assert.NotNil(suite.T(), errors.Cause(err))
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsNil_InCaseOfUpMigrationsToRun() {
	// Arrange
	// The following versions are from ../testdata.
	// We'll mark one of them as not migrated yet, meaning it needs
	// to be migrated up.
	migrations := map[int64]bool{
		1494538273: true,
		1494538317: false,
		1494538407: true,
	}
	files, err := file.ListFiles(filepath.Join("..", "testdata"), true)
	if err != nil {
		suite.FailNow(err.Error())
	}
	needsMigration := []file.File{
		*file.FindByVersion(1494538317, files),
	}
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	suite.driverMock.On("SelectMigrations", mock.AnythingOfType("*context.timerCtx")).Return(migrations, nil).Once()
	suite.driverMock.On("ApplyMigrations", mock.AnythingOfType("*context.timerCtx"), needsMigration, true).Return(nil).Once()
	suite.driverMock.On("Close").Once()

	// Act
	err = suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	assert.Nil(suite.T(), errors.Cause(err))
	assert.True(suite.T(), suite.output.Contains("1494538317_add-phone-number-to-users.up.sql"))
	assert.True(suite.T(), suite.output.Contains("seconds"))
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsNil_InCaseOfNoDownMigrationsToRun() {
	// Arrange
	// The following versions are from ../testdata.
	// We'll mark all of them as never been migrated, meaning
	// none of them need to be migrated down.
	migrations := map[int64]bool{
		1494538273: false,
		1494538317: false,
		1494538407: false,
	}
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	suite.driverMock.On("SelectMigrations", mock.AnythingOfType("*context.timerCtx")).Return(migrations, nil).Once()
	suite.driverMock.On("Close").Once()

	// Act
	err := suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", false)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	assert.Nil(suite.T(), errors.Cause(err))
	assert.True(suite.T(), suite.output.Contains("seconds"))
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsNil_InCaseOfDownMigrationsToRun() {
	// Arrange
	// The following versions are from ../testdata.
	// We'll mark one of them as migrated, meaning
	// it needs to be migrated down.
	migrations := map[int64]bool{
		1494538273: false,
		1494538317: false,
		1494538407: true,
	}
	files, err := file.ListFiles(filepath.Join("..", "testdata"), false)
	if err != nil {
		suite.FailNow(err.Error())
	}
	needsMigration := []file.File{
		*file.FindByVersion(1494538407, files),
	}
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	suite.driverMock.On("SelectMigrations", mock.AnythingOfType("*context.timerCtx")).Return(migrations, nil).Once()
	suite.driverMock.On("ApplyMigrations", mock.AnythingOfType("*context.timerCtx"), needsMigration, false).Return(nil).Once()
	suite.driverMock.On("Close").Once()

	// Act
	err = suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", false)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	assert.Nil(suite.T(), errors.Cause(err))
	assert.True(suite.T(), suite.output.Contains("1494538407_replace-user-phone-with-email.down.sql"))
	assert.True(suite.T(), suite.output.Contains("seconds"))
}

func (suite *MigratorTestSuite) Test_Create_ReturnsNil_InCaseOfSuccess() {
	// Arrange
	const path = "."

	// Act
	pair, err := suite.instance.Create("create_table_invoices", path)

	// Assert
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), pair)
	defer remove(filepath.Join(path, pair.Up.Base))
	defer remove(filepath.Join(path, pair.Down.Base))
	versionString := fmt.Sprintf("Version %d migration files created in %s", pair.Up.Version, path)
	assert.True(suite.T(), suite.output.Contains(versionString))
	assert.True(suite.T(), suite.output.Contains(pair.Up.Base))
	assert.True(suite.T(), suite.output.Contains(pair.Down.Base))
}

func remove(filename string) {
	if err := os.Remove(filename); err != nil {
		fmt.Println("removing file failed", err)
	}
}
