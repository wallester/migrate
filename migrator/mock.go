package migrator

import (
	"github.com/stretchr/testify/mock"
	"github.com/wallester/migrate/direction"
	"github.com/wallester/migrate/file"
)

// Mock is mock object for Migrator
type Mock struct {
	mock.Mock
}

// Migrate is a mock method
func (m *Mock) Migrate(path string, url string, dir direction.Direction, steps int, timeoutSeconds int) error {
	args := m.Called(path, url, dir, steps, timeoutSeconds)

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
