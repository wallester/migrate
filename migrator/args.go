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
	// Deprecated.  Use TimeoutDuration instead.
	TimeoutSeconds int
	URL            string
	Verbose        bool
}
