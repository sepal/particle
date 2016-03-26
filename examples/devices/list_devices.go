package main

import (
	"fmt"
	"github.com/sepal/particle"
	"github.com/sepal/particle/examples/common"
)

// Get the list of all the users devices.
func main() {
	token, err := common.GetToken()

	if err != nil {
		common.PrintError(err)
	}

	c := particle.NewClient(nil, token)

	devices, err := c.ListDevices()

	if err != nil {
		common.PrintError(err)
	}

	fmt.Println(devices)
}
