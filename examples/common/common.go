package common

import (
	"fmt"
	"github.com/mitchellh/colorstring"
	"os"
)

type UsageFunc func()

// PrintError exits the program with an error.
func PrintError(err error) {
	fmt.Println(colorstring.Color("[red]" + err.Error()))
	os.Exit(1)
}

// Exits and prints the app usage.
func UsageAndExit(message string, exitCode int, usage UsageFunc) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
