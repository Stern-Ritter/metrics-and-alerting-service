package pkg2

import (
	"os"
)

func main() {
	defer os.Exit(0)

	go func() {
		os.Exit(0)
	}()

	exit(0)

	os.Exit(0)
}

func exit(code int) {
	os.Exit(code)
}
