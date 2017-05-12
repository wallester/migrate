package driver

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wallester/migrate/file"
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

// SelectMigrations is a mock method
func (m *Mock) SelectMigrations(ctx context.Context) (map[int]bool, error) {
	args := m.Called(ctx)

	var result map[int]bool
	if args.Get(0) == nil {
		result = nil
	} else {
		result = args.Get(0).(map[int]bool)
	}

	return result, args.Error(1)
}

// ApplyMigrations is a mock method
func (m *Mock) ApplyMigrations(ctx context.Context, files []file.File, up bool) error {
	args := m.Called(ctx, files, up)

	return args.Error(0)
}

// Close is a mock method
func (m *Mock) Close() {
	m.Called()
}
