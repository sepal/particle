package main

import (
	"flag"
	"fmt"
	"github.com/sepal/particle/examples/common"
	"github.com/sepal/particle"
)

var token, deviceID, func_name, arg string

func main() {
	flag.StringVar(&token, "token", "", "Set the authentication token")
	flag.StringVar(&token, "t", "", "Set the authentication token (shorthand)")
	flag.StringVar(&deviceID, "device", "", "Set the device id")
	flag.StringVar(&deviceID, "d", "", "Set the device id (shorthand)")
	flag.StringVar(&func_name, "function", "", "Set the variable name to retrieve")
	flag.StringVar(&func_name, "f", "", "Set the variable name to retrieve (shorthand)")
	flag.StringVar(&arg, "argument", "", "Set a single argument to be passed to the function")
	flag.StringVar(&arg, "a", "", "Set a single argument to be passed to the function (shorthand)")

	flag.Usage = func() {
		fmt.Println("function -t token -d deviceID -f function_name [-a argument]")
		flag.PrintDefaults()
	}

	flag.Parse()

	if token == "" {
		common.UsageAndExit("Please set a token.", 0, flag.Usage)
	}

	if deviceID == "" {
		common.UsageAndExit("Please set a device ID.", 0, flag.Usage)
	}

	if func_name == "" {
		common.UsageAndExit("Please set a variable name to retrieve.", 0, flag.Usage)
	}


	c := particle.NewClient(nil, token)

	result, err := c.CallFunction(deviceID, func_name, arg)

	if err != nil {
		common.PrintError(err)
	}

	fmt.Printf("The function %v returned the value %v\n", func_name, result)
}