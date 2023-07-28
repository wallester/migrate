package direction

import (
	"github.com/mgutz/ansi"
)

// Direction is a boolean value that represents whether the database must be upgraded (Up) or downgraded (Down).
type Direction bool

const (
	Up   Direction = true
	Down Direction = false
)

// ToString returns the string representation of the direction.
func (d Direction) ToString() string {
	if d {
		return "up"
	}

	return "down"
}

func (d Direction) ToANSIColoredPrefix() string {
	if d {
		return ansi.Green + ">" + ansi.Reset
	}

	return ansi.Red + ">" + ansi.Reset
}
