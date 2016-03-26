package main

import (
	"github.com/sepal/particle/examples/common"
	"github.com/sepal/particle"
	"fmt"
)

// Get the list of all the users devices.
func main() {
	token, err := common.GetToken()

	if err != nil {
		common.PrintError(err)
	}

	c := particle.NewClient(nil ,token)

	devices, err := c.ListDevices()

	if err != nil {
		common.PrintError(err)
	}

	fmt.Println(devices)
}