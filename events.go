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

type EventChannel chan Event

type Event struct {
	Name        string
	Data        string
	TTL         string
	PublishedAt time.Time `json:"published_at"`
}

type EventListener struct {
	OutputChan EventChannel
	response   *http.Response
	running    bool
}

func (c *Client) NewEventListener(name, deviceID string) (*EventListener, error) {
	e := &EventListener{}

	if e.OutputChan == nil {
		e.OutputChan = make(chan Event)
	}

	if e.response == nil {
		var endPoint string
		if deviceID != "" {
			endPoint = deviceURL + "/event"
		} else {
			endPoint = eventURL
		}

		if name != "" {
			endPoint += "/" + name
		}

		resp, err := c.Get(endPoint, nil)

		if err != nil {
			return nil, err
		}

		err = CheckResponse(resp)

		if err != nil {
			return nil, err
		}

		e.response = resp
	}

	return e, nil
}

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
			b := buf.Bytes()
			err := json.Unmarshal(b, &ev)

			if err == nil {
				e.OutputChan <- ev
			}
			buf.Reset()
			ev = Event{}
			// todo: Error handling for if json decoding failed.
		}

	}

	return nil
}

func (e *EventListener) Close() {
	close(e.OutputChan)
	e.running = false
}
