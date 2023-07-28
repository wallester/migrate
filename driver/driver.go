package driver

import (
	"context"

	"github.com/wallester/migrate/direction"
	"github.com/wallester/migrate/file"
	"github.com/wallester/migrate/version"
)

// Driver represents database driver interface.
type IDriver interface {
	Open(ctx context.Context, url string) error
	CreateMigrationsTable(ctx context.Context) error
	SelectAllMigrations(ctx context.Context) (version.Versions, error)
	Migrate(ctx context.Context, f file.File, d direction.Direction) error
	Close() error
}
