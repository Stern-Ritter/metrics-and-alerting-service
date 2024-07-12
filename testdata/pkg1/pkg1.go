package main

import (
	"os"
)

func main() {
	defer os.Exit(0) // want `using os.Exit call in main func at line 8, column 8, full call expr: os.Exit\(0\)`

	go func() {
		os.Exit(0) // want `using os.Exit call in main func at line 11, column 3, full call expr: os.Exit\(0\)`
	}()

	exit(0)

	os.Exit(0) // want `using os.Exit call in main func at line 16, column 2, full call expr: os.Exit\(0\)`
}

func exit(code int) {
	os.Exit(code)
}
