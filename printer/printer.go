package printer

import (
	"fmt"
)

// Printer prints
type Printer interface {
	Println(a ...interface{})
}

// New returns new instance
func New() Printer {
	return &printer{}
}

type printer struct{}

func (p *printer) Println(a ...interface{}) {
	fmt.Println(a...)
}
