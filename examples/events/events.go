package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/colorstring"
	"github.com/sepal/particle"
	"github.com/sepal/particle/examples/common"
	"os"
	"os/signal"
	"syscall"
)

var token, event, deviceID string

func main() {
	flag.StringVar(&token, "token", "", "Set the authentication token")
	flag.StringVar(&token, "t", "", "Set the authentication token (shorthand)")
	flag.StringVar(&event, "event", "", "Event name to subcribe to")
	flag.StringVar(&event, "e", "", "Event name to subcribe to (shorthand)")
	flag.StringVar(&deviceID, "device", "", "Set the device id")
	flag.StringVar(&deviceID, "d", "", "Set the device id (shorthand)")

	flag.Usage = func() {
		fmt.Println("events -t token [-e event]")
		flag.PrintDefaults()
	}

	flag.Parse()

	if token == "" {
		common.UsageAndExit("Please set a token.", 0, flag.Usage)
	}

	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	signal.Notify(osChan, syscall.SIGTERM)

	c := particle.NewClient(nil, token)

	e, err := c.NewEventListener(event, deviceID)

	if err != nil {
		common.PrintError(err)
	}

	fmt.Println(colorstring.Color("[green]Starting listening."))
	go e.Listen()

	go func() {
		<-osChan
		e.Close()
		fmt.Println(colorstring.Color("[green]Closing event listener."))
	}()

	for event := range e.OutputChan {
		fmt.Printf("New event %v, with data: %v\n", event.Name, event.Data)
	}

	fmt.Println(colorstring.Color("[green]Finished listening."))
}
