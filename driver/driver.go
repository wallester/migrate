package driver

import (
	"context"

	"github.com/wallester/migrate/file"
	"github.com/wallester/migrate/version"
)

// Driver represents database driver
type IDriver interface {
	Open(url string) error
	CreateMigrationsTable(ctx context.Context) error
	SelectAllMigrations(ctx context.Context) (version.Versions, error)
	Migrate(ctx context.Context, f file.File, up bool) error
	Close() error
}
