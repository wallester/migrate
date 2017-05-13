package driver

import (
	"context"

	"github.com/wallester/migrate/file"
)

// Driver represents database driver
type Driver interface {
	Open(url string) error
	CreateMigrationsTable(ctx context.Context) error
	SelectAllMigrations(ctx context.Context) (map[int64]bool, error)
	ApplyMigrations(ctx context.Context, files []file.File, up bool) error
	Close()
}
