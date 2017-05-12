package migrator

import (
	"github.com/stretchr/testify/mock"
)

// Mock is mock object for Migrator
type Mock struct {
	mock.Mock
}

// Up is a mock method
func (m *Mock) Up(path string, url string) error {
	args := m.Called(path, url)

	return args.Error(0)
}

// Down is a mock method
func (m *Mock) Down(path string, url string) error {
	args := m.Called(path, url)

	return args.Error(0)
}

// Create is a mock method
func (m *Mock) Create(name string, path string) error {
	args := m.Called(name, path)

	return args.Error(0)
}
