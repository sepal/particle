package particle

import (
	"time"
	"net/http"
)

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

func (c *Client) NewEventListener(eventName string) (*EventListener, error)  {
	return nil, nil
}