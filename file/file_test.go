package file

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ListFiles_ReturnsUpMigrationFiles_InCaseOfSuccess(t *testing.T) {
	// Act
	files, err := ListFiles(filepath.Join("..", "testdata"), true)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 3, len(files))

	assert.Equal(t, "1494538273_create-table-users.up.sql", files[0].Base)
	assert.Equal(t, 1494538273, files[0].Version)
	assert.NotEmpty(t, files[0].SQL)

	assert.Equal(t, "1494538317_add-phone-number-to-users.up.sql", files[1].Base)
	assert.Equal(t, 1494538317, files[1].Version)
	assert.NotEmpty(t, files[1].SQL)

	assert.Equal(t, "1494538407_replace-user-phone-with-email.up.sql", files[2].Base)
	assert.Equal(t, 1494538407, files[2].Version)
	assert.NotEmpty(t, files[2].SQL)
}

func Test_ListFiles_ReturnsDownMigrationFiles_InCaseOfSuccess(t *testing.T) {
	// Act
	files, err := ListFiles(filepath.Join("..", "testdata"), false)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 3, len(files))

	assert.Equal(t, "1494538407_replace-user-phone-with-email.down.sql", files[0].Base)
	assert.Equal(t, 1494538407, files[0].Version)
	assert.NotEmpty(t, files[0].SQL)

	assert.Equal(t, "1494538317_add-phone-number-to-users.down.sql", files[1].Base)
	assert.Equal(t, 1494538317, files[1].Version)
	assert.NotEmpty(t, files[1].SQL)

	assert.Equal(t, "1494538273_create-table-users.down.sql", files[2].Base)
	assert.Equal(t, 1494538273, files[2].Version)
	assert.NotEmpty(t, files[2].SQL)
}
