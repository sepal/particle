package main

import (
	"fmt"
	"github.com/mckee/particle"
)

const URL string = "https://api.particle.io/v1/devices/events/"

const token string = "2f556a89ad71889f2195f85bc787db46aa5f6804"
const particleDevice string = "55ff66065075555342421787"

func main() {
	event_channel := particle.Subscribe("/weather", token)
	fmt.Print(event_channel)
}
