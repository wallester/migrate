package driver

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wallester/migrate/file"
	"github.com/wallester/migrate/version"
)

// Mock is mock object for Driver
type Mock struct {
	mock.Mock
}

// Open is a mock method
func (m *Mock) Open(url string) error {
	args := m.Called(url)

	return args.Error(0)
}

// CreateMigrationsTable is a mock method
func (m *Mock) CreateMigrationsTable(ctx context.Context) error {
	args := m.Called(ctx)

	return args.Error(0)
}

// SelectAllMigrations is a mock method
func (m *Mock) SelectAllMigrations(ctx context.Context) (version.Versions, error) {
	args := m.Called(ctx)

	if args.Get(0) != nil {
		return args.Get(0).(version.Versions), args.Error(1)
	}

	return nil, args.Error(1)
}

// ApplyMigrations is a mock method
func (m *Mock) ApplyMigrations(ctx context.Context, files []file.File, up bool) error {
	args := m.Called(ctx, files, up)

	return args.Error(0)
}

// Close is a mock method
func (m *Mock) Close() error {
	args := m.Called()

	return args.Error(0)
}
