package main

import (
	"fmt"
	"os"

	"github.com/wallester/migrate/app"
)

func main() {
	if err := app.New().Run(os.Args); err != nil {
		//nolint:forbidigo
		fmt.Println(err)
		os.Exit(1)
	}
}
