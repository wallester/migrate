package printer

import (
	"fmt"
	"strings"
)

// Recorder records printer output
type Recorder struct {
	Output []string
}

var _ IPrinter = (*Recorder)(nil)

// Println remembers the printed values
func (r *Recorder) Println(a ...interface{}) {
	line := make([]string, 0, len(a))
	for _, v := range a {
		line = append(line, fmt.Sprintf("%v", v))
	}

	r.Output = append(r.Output, strings.Join(line, " "))
}

// String returns output as string
func (r *Recorder) String() string {
	return strings.Join(r.Output, "\n")
}

// Contains returns true if output contains given text
func (r *Recorder) Contains(text string) bool {
	return strings.Contains(r.String(), text)
}
