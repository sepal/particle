package particle

import (
	"time"
	"net/http"
)

const eventURL  = "/v1/events"

type EventChannel chan Event

type Event struct {
	Name string
	Data string
	TTL int
	PublishedAt time.Time `json:"published_at"`
}

type EventListener struct {
	OutputChan EventChannel
	Response *http.Response
	running bool
}

func (c *Client) NewEventListener(name string) (*EventListener, error)  {
	e := &EventListener{}
	
	if e.OutputChan == nil {
		e.OutputChan = make(chan Event)
	}
	
	if e.Response == nil {
		resp, err := c.Get(eventURL + "/" + name, nil)

		if err != nil {
			return nil, err
		}

		err = CheckResponse(resp)

		if  err != nil {
			return nil, err
		}

		e.Response = resp
	}
	
	return e, nil
}