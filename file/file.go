package file

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/juju/errors"
)

var filePrefix = map[bool]string{
	true:  "up",
	false: "down",
}

// File represents a migration file
type File struct {
	Base    string
	Version int
	SQL     string
}

// ByBase implements sort.Interface for []MigrationFile based on
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

// ListFiles lists migration files on a given path
func ListFiles(path string, up bool) ([]File, error) {
	files, err := filepath.Glob(filepath.Join(path, "*_*."+filePrefix[up]+".sql"))
	if err != nil {
		return nil, errors.Annotate(err, "getting migration files failed")
	}

	var migrations []File
	for _, file := range files {
		base := filepath.Base(file)

		version, err := strconv.Atoi(strings.Split(base, "_")[0])
		if err != nil {
			return nil, errors.Annotate(err, "parsing version failed")
		}

		b, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, errors.Annotate(err, "reading migration file failed")
		}

		migrations = append(migrations, File{
			Base:    base,
			Version: version,
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
