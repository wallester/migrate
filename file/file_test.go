package file

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wallester/migrate/direction"
)

// nolint: dupl
func Test_ListFiles_ReturnsUpMigrationFiles_InCaseOfSuccess(t *testing.T) {
	// Act
	files, err := ListFiles(filepath.Join("..", "testdata"), direction.Up)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 3, len(files))

	assert.Equal(t, "1494538273_create_table_users.up.sql", files[0].Base)
	assert.Equal(t, int64(1494538273), files[0].Version)
	assert.NotEmpty(t, files[0].SQL)

	assert.Equal(t, "1494538317_add_phone_number_to_users.up.sql", files[1].Base)
	assert.Equal(t, int64(1494538317), files[1].Version)
	assert.NotEmpty(t, files[1].SQL)

	assert.Equal(t, "1494538407_replace_user_phone_with_email.up.sql", files[2].Base)
	assert.Equal(t, int64(1494538407), files[2].Version)
	assert.NotEmpty(t, files[2].SQL)
}

// nolint: dupl
func Test_ListFiles_ReturnsDownMigrationFiles_InCaseOfSuccess(t *testing.T) {
	// Act
	files, err := ListFiles(filepath.Join("..", "testdata"), direction.Down)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 3, len(files))

	assert.Equal(t, "1494538407_replace_user_phone_with_email.down.sql", files[0].Base)
	assert.Equal(t, int64(1494538407), files[0].Version)
	assert.NotEmpty(t, files[0].SQL)

	assert.Equal(t, "1494538317_add_phone_number_to_users.down.sql", files[1].Base)
	assert.Equal(t, int64(1494538317), files[1].Version)
	assert.NotEmpty(t, files[1].SQL)

	assert.Equal(t, "1494538273_create_table_users.down.sql", files[2].Base)
	assert.Equal(t, int64(1494538273), files[2].Version)
	assert.NotEmpty(t, files[2].SQL)
}
