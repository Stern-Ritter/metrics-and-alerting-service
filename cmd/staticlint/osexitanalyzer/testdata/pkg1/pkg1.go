package main

import (
	"os"
)

func main() {
	defer os.Exit(0) // want "using os.Exit in main func"

	go func() {
		os.Exit(0) // want "using os.Exit in main func"
	}()

	exit(0)

	os.Exit(0) // want "using os.Exit in main func"
}

func exit(code int) {
	os.Exit(code)
}
