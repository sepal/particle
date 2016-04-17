package particle

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const eventURL = "/v1/events"

var eventNameLabel = []byte("event:")
var eventDataLabel = []byte("data:")

// EventChannel are the channels where the Events outputted to.
type EventChannel chan Event

// ErrorChannel receive any errors that occurred while an EventListener is listening.
type ErrorChannel chan error

// Event represents a single event from the particle cloud api.
type Event struct {
	Name        string
	Data        string
	TTL         string
	PublishedAt time.Time `json:"published_at"`
}

// EventListener listens to events from the particle cloud api and outputs them over the OutputChan channel.
type EventListener struct {
	OutputChan EventChannel
	ErrorChan  ErrorChannel
	response   *http.Response
	running    bool
}

// newEventListener creates a new EventListener + its channels but doesn't connects to the server
func newEventListener() *EventListener {
	e := &EventListener{}

	e.OutputChan = make(chan Event)
	e.ErrorChan = make(chan error)

	return e
}

// connectEventListener connects the given EventListener to the given endPoint.
func (c *Client) connectEventListener(endPoint string, e *EventListener) error {
	resp, err := c.get(endPoint, nil)

	if err != nil {
		return err
	}

	err = CheckResponse(resp)

	if err != nil {
		return err
	}

	e.response = resp
	return nil
}

// NewEventListener creates a new EventListener for the given event and device id. Both parameters are optional, the
// event listener will then listen all events. The function will also connect to the server.
func (c *Client) NewEventListener(name string) (*EventListener, error) {
	e := newEventListener()

	endPoint := eventURL

	if name != "" {
		endPoint += "/" + name
	}

	err := c.connectEventListener(endPoint, e)

	if err != nil {
		return nil, err
	}

	return e, nil
}

// NewPrivateEventListener creates a new EventListener, which subscribes events for devices of the the tokens account.
func (c *Client) NewPrivateEventListener(name string) (*EventListener, error) {
	e := newEventListener()

	endPoint := deviceURL + "/events"

	if name != "" {
		endPoint += "/" + name
	}

	err := c.connectEventListener(endPoint, e)

	if err != nil {
		return nil, err
	}

	return e, nil
}

// NewEventListener creates a new EventListener for this device for the given event name. If the name is omitted then
// the function will subscribe to all events of this device.
func (d *Device) NewEventListener(name string) (*EventListener, error) {
	e := newEventListener()

	if d.ID == "" {
		return nil, fmt.Errorf("Device %v has no id", d)
	}

	endPoint := deviceURL + "/" + d.ID + "/events"

	if name != "" {
		endPoint += "/" + name
	}

	err := d.client.connectEventListener(endPoint, e)

	if err != nil {
		return nil, err
	}

	return e, nil
}

// Listen starts reading events from the cloud API.
func (e *EventListener) Listen() error {
	ev := Event{}
	reader := bufio.NewReader(e.response.Body)
	var buf bytes.Buffer

	e.running = true

	for e.running {
		line, err := reader.ReadBytes('\n')

		if err != nil {
			return fmt.Errorf("Error while reading line: %v", err)
		}

		switch {
		// todo: check for :ok
		case bytes.HasPrefix(line, eventNameLabel):
			ev.Name = string(line[len(eventNameLabel):])
			ev.Name = strings.TrimSpace(ev.Name)
		case bytes.HasPrefix(line, eventDataLabel):
			buf.Write(line[len(eventDataLabel):])
		case bytes.Equal(line, []byte("\n")):
			if buf.Len() > 0 {
				b := buf.Bytes()
				err := json.Unmarshal(b, &ev)

				if err == nil {
					e.OutputChan <- ev
				} else {
					e.ErrorChan <- err
				}

				buf.Reset()
				ev = Event{}
			}
		}
	}

	return nil
}

// Close closes the EventListeners channel and stops the listening loop.
func (e *EventListener) Close() {
	close(e.OutputChan)
	e.running = false
}
