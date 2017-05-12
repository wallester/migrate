package printer

import (
	"fmt"
)

// Recorder records printer output
type Recorder struct {
	output []string
}

// Println remembers the printed values
func (r *Recorder) Println(a ...interface{}) {
	for _, v := range a {
		r.output = append(r.output, fmt.Sprintf("%v", v))
	}
}
