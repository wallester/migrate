package driver

import (
	"context"

	"github.com/wallester/migrate/direction"
	"github.com/wallester/migrate/file"
	"github.com/wallester/migrate/version"
)

// Driver represents database driver
type Driver interface {
	Open(url string) error
	CreateMigrationsTable(ctx context.Context) error
	SelectAllMigrations(ctx context.Context) (version.Versions, error)
	ApplyMigrations(ctx context.Context, files []file.File, dir direction.Direction) error
	Close() error
}
