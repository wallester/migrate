package printer

import (
	"fmt"
)

// Printer prints
type IPrinter interface {
	Println(a ...interface{})
}

type Printer struct{}

var _ IPrinter = (*Printer)(nil)

// New returns new instance
func New() *Printer {
	return &Printer{}
}

func (p *Printer) Println(a ...interface{}) {
	//nolint:forbidigo
	fmt.Println(a...)
}
