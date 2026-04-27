package main

import (
	"os"

	"sysi/internal/sysiapp"
)

func main() {
	os.Exit(sysiapp.New(sysiapp.Options{}).Run(os.Args[1:]))
}
