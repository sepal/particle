package main

import (
	"flag"
	"fmt"
	"github.com/sepal/particle"
	"github.com/sepal/particle/examples/common"
)

var token string

// Get the list of all the users devices.
func main() {
	flag.StringVar(&token, "token", "", "Set the authentication token")
	flag.StringVar(&token, "t", "", "Set the authentication token (shorthand)")

	flag.Usage = func() {
		fmt.Println("list_devices -t [token]")
		flag.PrintDefaults()
	}

	flag.Parse()

	if token == "" {
		common.UsageAndExit("Please set a token.", 0, flag.Usage)
	}

	c := particle.NewClient(nil, token)

	devices, err := c.ListDevices()

	if err != nil {
		common.PrintError(err)
	}

	fmt.Println(devices)
}
