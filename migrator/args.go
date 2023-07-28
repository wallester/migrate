package migrator

import (
	"time"

	"github.com/wallester/migrate/direction"
)

type Args struct {
	DBConnectionTimeoutDuration time.Duration
	Direction                   direction.Direction
	NoVerify                    bool
	Path                        string
	Steps                       int
	TimeoutDuration             time.Duration
	URL                         string
	Verbose                     bool
}
