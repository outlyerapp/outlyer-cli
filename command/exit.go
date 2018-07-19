package command

import (
	"fmt"
	"os"
)

const (
	// ExitError represents a general error (http://tldp.org/LDP/abs/html/exitcodes.html)
	ExitError = 1
	// ExitBadArgs represents invalid arguments error
	ExitBadArgs = 128
)

// ExitWithSuccess prints a message to stdout and exits with code 0
func ExitWithSuccess(msg string) {
	fmt.Fprintln(os.Stdout, msg)
	os.Exit(0)
}

// ExitWithError prints an error message to stderr and exits with the specified code
func ExitWithError(code int, err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(code)
}
