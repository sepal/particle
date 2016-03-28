package main

import (
	"fmt"
	"github.com/sepal/particle"
	"github.com/sepal/particle/examples/common"
	"flag"
	"reflect"
)

var token, deviceID, variable, varType string

// Get the list of all the users devices.
func main() {
	flag.StringVar(&token, "token", "", "Set the authentication token")
	flag.StringVar(&token, "t", "", "Set the authentication token (shorthand)")
	flag.StringVar(&deviceID, "device", "", "Set the device id")
	flag.StringVar(&deviceID, "d", "", "Set the device id (shorthand)")
	flag.StringVar(&variable, "variable", "", "Set the variable name to retrieve")
	flag.StringVar(&variable, "v", "", "Set the variable name to retrieve (shorthand)")
	flag.StringVar(&varType, "type", "string", "Set the expected type of the variable value.")

	flag.Usage = func() {
		fmt.Println("variables -t [token] -d [deviceID] -v [variable_name]")
		flag.PrintDefaults()
	}

	flag.Parse()

	if token == "" {
		common.UsageAndExit("Please set a token.", 0, flag.Usage)
	}

	if deviceID == "" {
		common.UsageAndExit("Please set a device ID.", 0, flag.Usage)
	}

	if variable == "" {
		common.UsageAndExit("Please set a variable name to retrieve.", 0, flag.Usage)
	}

	c := particle.NewClient(nil, token)

	switch varType {
	case "string":
		value, err := c.VariableString(deviceID, variable)
		if err != nil {
			common.PrintError(err)
		}

		fmt.Printf("Value for '%v' from device '%v' is '%v' with the type '%v'.\n", variable, deviceID, value,
			reflect.TypeOf(value))
	case "int":
		value, err := c.VariableInt(deviceID, variable)
		if err != nil {
			common.PrintError(err)
		}

		fmt.Printf("Value for '%v' from device '%v' is '%v with the type '%v'.\n", variable, deviceID, value,
			reflect.TypeOf(value))
	case "float":
		value, err := c.VariableFloat(deviceID, variable)
		if err != nil {
			common.PrintError(err)
		}

		fmt.Printf("Value for '%v' from device '%v' is '%v' with the type '%v'.\n", variable, deviceID, value,
			reflect.TypeOf(value))
	default:
		msg := fmt.Sprintf("The passed type '%v' is not supported. Please use 'string', 'int' or 'float'",
			varType);
		common.UsageAndExit(msg, 1, flag.Usage)
	}
}
