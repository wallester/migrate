package main

import (
	"fmt"
	"os"

	"github.com/wallester/migrate/app"
)

func main() {
	if err := app.New().Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
