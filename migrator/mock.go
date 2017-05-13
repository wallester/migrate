package migrator

import (
	"github.com/stretchr/testify/mock"
	"github.com/wallester/migrate/file"
)

// Mock is mock object for Migrator
type Mock struct {
	mock.Mock
}

// MigrateAll is a mock method
func (m *Mock) MigrateAll(path string, url string, up bool) error {
	args := m.Called(path, url, up)

	return args.Error(0)
}

// Create is a mock method
func (m *Mock) Create(name string, path string) (*file.Pair, error) {
	args := m.Called(name, path)

	if args.Get(0) != nil {
		return args.Get(0).(*file.Pair), args.Error(1)
	}

	return nil, args.Error(1)
}
