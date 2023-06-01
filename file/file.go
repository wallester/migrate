package file

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/juju/errors"
	"github.com/wallester/migrate/direction"
)

// File represents a migration file
type File struct {
	Base    string
	Version int64
	SQL     string
}

// Create creates a new file in the given path
func (f File) Create(path string) error {
	if err := os.WriteFile(filepath.Join(path, f.Base), nil, 0o600); err != nil {
		return errors.Annotate(err, "writing migration file failed")
	}

	return nil
}

// Pair is a pair of migration files; up and down
type Pair struct {
	Up   File
	Down File
}

// ByBase implements sort.Interface for []File based on
// the Base field.
type ByBase []File

func (a ByBase) Len() int {
	return len(a)
}

func (a ByBase) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByBase) Less(i, j int) bool {
	return a[i].Base < a[j].Base
}

// FindByVersion finds a file from list by version
func FindByVersion(version int64, files []File) *File {
	for _, file := range files {
		if file.Version == version {
			return &file
		}
	}

	return nil
}

// ListFiles lists migration files on a given path
func ListFiles(path string, up bool) ([]File, error) {
	files, err := filepath.Glob(filepath.Join(path, "*_*."+fileSuffix[up]+".sql"))
	if err != nil {
		return nil, errors.Annotate(err, "getting migration files failed")
	}

	migrations := make([]File, 0, len(files))
	for _, file := range files {
		base := filepath.Base(file)

		version, err := version(base)
		if err != nil {
			return nil, errors.Annotatef(err, "getting version of %s migration failed", base)
		}

		b, err := os.ReadFile(filepath.Clean(file))
		if err != nil {
			return nil, errors.Annotate(err, "reading migration file failed")
		}

		migrations = append(migrations, File{
			Base:    base,
			Version: *version,
			SQL:     string(b),
		})
	}

	if up {
		sort.Sort(ByBase(migrations))
	} else {
		sort.Sort(sort.Reverse(ByBase(migrations)))
	}

	return migrations, nil
}

// private

// version returns version of migration file
func version(base string) (*int64, error) {
	version, err := strconv.ParseInt(strings.Split(base, "_")[0], 10, 64)
	if err != nil {
		return nil, errors.Annotate(err, "parsing version failed")
	}

	return &version, nil
}

var fileSuffix = map[bool]string{
	direction.Up:   "up",
	direction.Down: "down",
}
