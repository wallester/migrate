package migrator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/wallester/migrate/driver"
	"github.com/wallester/migrate/file"
	"github.com/wallester/migrate/printer"
	"github.com/wallester/migrate/version"
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
	suite.NotNil(instance)
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsNil_InCaseOfNoUpMigrationsToRun() {
	// Arrange
	// The following versions are from ../testdata.
	// We'll mark all of them as already migrated, meaning
	// no up migrations need to run.
	var exists struct{}
	migrations := version.Versions{
		1494538273: exists,
		1494538317: exists,
		1494538407: exists,
	}
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	suite.driverMock.On("SelectAllMigrations", mock.AnythingOfType("*context.timerCtx")).Return(migrations, nil).Once()
	suite.driverMock.On("Close").Return(nil).Once()

	// Act
	err := suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true, 0, 1)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	suite.Nil(errors.Cause(err))
	suite.True(suite.output.Contains("seconds"))
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsError_InCaseOfDriverOpenError() {
	// Arrange
	suite.driverMock.On("Open", "connectionurl").Return(suite.expectedErr).Once()

	// Act
	err := suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true, 0, 1)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	suite.NotNil(err)
	suite.EqualError(err, "opening database connection failed: failure")
	suite.Empty(suite.output.String())
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsError_InCaseOfDriverCreateMigrationsTableError() {
	// Arrange
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(suite.expectedErr).Once()
	suite.driverMock.On("Close").Return(nil).Once()

	// Act
	err := suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true, 0, 1)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	suite.NotNil(err)
	suite.EqualError(err, "migrating failed: creating migrations table failed: failure")
	suite.Empty(suite.output.String())
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsErr_InCaseOfDriverSelectMigrationsError() {
	// Arrange
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	suite.driverMock.On("SelectAllMigrations", mock.AnythingOfType("*context.timerCtx")).Return(nil, suite.expectedErr).Once()
	suite.driverMock.On("Close").Return(nil).Once()

	// Act
	err := suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true, 0, 1)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	suite.NotNil(err)
	suite.EqualError(err, "migrating failed: selecting existing migrations failed: failure")
	suite.Empty(suite.output.String())
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsError_InCaseOfDriverApplyMigrationsError() {
	// Arrange
	// The following versions are from ../testdata.
	// We'll mark one of them as not migrated yet, meaning it needs
	// to be migrated up.
	var exists struct{}
	migrations := version.Versions{
		1494538273: exists,
		1494538317: exists,
	}
	files, err := file.ListFiles(filepath.Join("..", "testdata"), true)
	if err != nil {
		suite.FailNow(err.Error())
	}
	needsMigration := []file.File{
		*file.FindByVersion(1494538407, files),
	}
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	suite.driverMock.On("SelectAllMigrations", mock.AnythingOfType("*context.timerCtx")).Return(migrations, nil).Once()
	suite.driverMock.On("ApplyMigrations", mock.AnythingOfType("*context.timerCtx"), needsMigration, true).Return(suite.expectedErr).Once()
	suite.driverMock.On("Close").Return(nil).Once()

	// Act
	err = suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true, 0, 1)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	suite.NotNil(errors.Cause(err))
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsNil_InCaseOfUpMigrationsToRun() {
	// Arrange
	// The following versions are from ../testdata.
	// We'll mark one of them as not migrated yet, meaning it needs
	// to be migrated up.
	var exists struct{}
	migrations := version.Versions{
		1494538273: exists,
		1494538317: exists,
	}
	files, err := file.ListFiles(filepath.Join("..", "testdata"), true)
	if err != nil {
		suite.FailNow(err.Error())
	}
	needsMigration := []file.File{
		*file.FindByVersion(1494538407, files),
	}
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	suite.driverMock.On("SelectAllMigrations", mock.AnythingOfType("*context.timerCtx")).Return(migrations, nil).Once()
	suite.driverMock.On("ApplyMigrations", mock.AnythingOfType("*context.timerCtx"), needsMigration, true).Return(nil).Once()
	suite.driverMock.On("Close").Return(nil).Once()

	// Act
	err = suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true, 0, 1)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	suite.Nil(errors.Cause(err))
	suite.True(suite.output.Contains("1494538407_replace_user_phone_with_email.up.sql"))
	suite.True(suite.output.Contains("seconds"))
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsNil_InCaseOfNoDownMigrationsToRun() {
	// Arrange
	// The following versions are from ../testdata.
	// We'll mark all of them as never been migrated, meaning
	// none of them need to be migrated down.
	migrations := make(version.Versions)
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	suite.driverMock.On("SelectAllMigrations", mock.AnythingOfType("*context.timerCtx")).Return(migrations, nil).Once()
	suite.driverMock.On("Close").Return(nil).Once()

	// Act
	err := suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", false, 0, 1)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	suite.Nil(errors.Cause(err))
	suite.True(suite.output.Contains("seconds"))
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsNil_InCaseOfDownMigrationsToRun() {
	// Arrange
	// The following versions are from ../testdata.
	// We'll mark one of them as migrated, meaning
	// it needs to be migrated down.
	var exists struct{}
	migrations := version.Versions{
		1494538407: exists,
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
	suite.driverMock.On("SelectAllMigrations", mock.AnythingOfType("*context.timerCtx")).Return(migrations, nil).Once()
	suite.driverMock.On("ApplyMigrations", mock.AnythingOfType("*context.timerCtx"), needsMigration, false).Return(nil).Once()
	suite.driverMock.On("Close").Return(nil).Once()

	// Act
	err = suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", false, 0, 1)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	suite.Nil(errors.Cause(err))
	suite.True(suite.output.Contains("1494538407_replace_user_phone_with_email.down.sql"))
	suite.True(suite.output.Contains("seconds"))
}

func (suite *MigratorTestSuite) Test_Create_ReturnsNil_InCaseOfSuccess() {
	// Arrange
	const path = "."

	// Act
	pair, err := suite.instance.Create("create_table_invoices", path)

	// Assert
	suite.Nil(err)
	suite.NotNil(pair)
	defer remove(filepath.Join(path, pair.Up.Base))
	defer remove(filepath.Join(path, pair.Down.Base))
	versionString := fmt.Sprintf("Version %d migration files created in %s", pair.Up.Version, path)
	suite.True(suite.output.Contains(versionString))
	suite.True(suite.output.Contains(pair.Up.Base))
	suite.True(suite.output.Contains(pair.Down.Base))
}

func remove(filename string) {
	if err := os.Remove(filename); err != nil {
		fmt.Println("removing file failed", err)
	}
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsNil_InCaseOfOneUpMigrationToRun() {
	// Arrange
	// The following versions are from ../testdata.
	migrations := make(version.Versions)
	files, err := file.ListFiles(filepath.Join("..", "testdata"), true)
	if err != nil {
		suite.FailNow(err.Error())
	}
	needsMigration := []file.File{
		*file.FindByVersion(1494538273, files),
	}
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	suite.driverMock.On("SelectAllMigrations", mock.AnythingOfType("*context.timerCtx")).Return(migrations, nil).Once()
	suite.driverMock.On("ApplyMigrations", mock.AnythingOfType("*context.timerCtx"), needsMigration, true).Return(nil).Once()
	suite.driverMock.On("Close").Return(nil).Once()

	// Act
	err = suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true, 1, 1)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	suite.Nil(errors.Cause(err))
	suite.True(suite.output.Contains("1494538273_create_table_users.up.sql"))
	suite.True(suite.output.Contains("seconds"))
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsNil_InCaseOfOneDownMigrationToRun() {
	// Arrange
	// The following versions are from ../testdata.
	var exists struct{}
	migrations := version.Versions{
		1494538273: exists,
		1494538317: exists,
		1494538407: exists,
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
	suite.driverMock.On("SelectAllMigrations", mock.AnythingOfType("*context.timerCtx")).Return(migrations, nil).Once()
	suite.driverMock.On("ApplyMigrations", mock.AnythingOfType("*context.timerCtx"), needsMigration, false).Return(nil).Once()
	suite.driverMock.On("Close").Return(nil).Once()

	// Act
	err = suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", false, 1, 1)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	suite.Nil(errors.Cause(err))
	suite.True(suite.output.Contains("1494538407_replace_user_phone_with_email.down.sql"))
	suite.True(suite.output.Contains("seconds"))
}

func (suite *MigratorTestSuite) Test_Migrate_ReturnsError_InCaseOfUpMigrationOlderThanAlreadyMigratedOne() {
	// Arrange
	// The following versions are from ../testdata.
	// We'll mark one of them as not migrated yet, meaning it needs
	// to be migrated up.
	var exists struct{}
	migrations := version.Versions{
		1494538273: exists,
		1494538407: exists,
	}
	suite.driverMock.On("Open", "connectionurl").Return(nil).Once()
	suite.driverMock.On("CreateMigrationsTable", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	suite.driverMock.On("SelectAllMigrations", mock.AnythingOfType("*context.timerCtx")).Return(migrations, nil).Once()
	suite.driverMock.On("Close").Return(nil).Once()

	// Act
	err := suite.instance.Migrate(filepath.Join("..", "testdata"), "connectionurl", true, 0, 1)

	// Assert
	suite.driverMock.AssertExpectations(suite.T())
	suite.NotNil(errors.Cause(err))
	suite.EqualError(errors.Cause(err), "cannot migrate up 1494538317_add_phone_number_to_users.up.sql, because it's older than already migrated version 1494538407")
}
