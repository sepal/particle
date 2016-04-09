package particle

import (
	"testing"
	"net/http"
)

func TestClient_NewEventListener(t *testing.T) {
	setup()
	defer teardown()

	eventName := "some_event"

	mux.HandleFunc(eventURL + "/" + eventName, func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; r.Method != m {
			t.Errorf("Wrong request method %v, expected %v", r.Method, m)
		}
	})

	e, err := client.NewEventListener(eventName)

	if err != nil {
		t.Fatalf("Error while creating EventListener: %v", err)
	}

	if e.OutputChan == nil {
		t.Errorf("NewEventListener didn't create an output channel.")
	}

	if e.Response == nil {
		t.Errorf("The EventListeners response is nil.")
	}

	if e.running {
		t.Errorf("EventListener is already listening although .listen wasn't called.")
	}
}
