package common

import (
	"os"
	"fmt"
	"errors"
	"github.com/mitchellh/colorstring"
)

// GetToken returns the token parsed from either the environment variables or the arguments.
// The arguments have priority over the environment variable.
func GetToken() (string, error) {
	var token string
	var err error

	if len(os.Args) == 3 {
		if os.Args[1] == "-t" {
			token = os.Args[2]
		}
	} else {
		token = os.Getenv("TOKEN")
	}

	if token == "" {
		err = errors.New("You have to provide a token either with then TOKEN env or the -t argument.")
	}

	return token, err
}

// Exit with an error.
func PrintError(err error) {
	fmt.Println(colorstring.Color("[red]" + err.Error()))
	os.Exit(1)
}