package particle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestClient_NewEventListener(t *testing.T) {
	setup()
	defer teardown()

	eventName := "some_event"

	mux.HandleFunc(eventURL+"/"+eventName, func(w http.ResponseWriter, r *http.Request) {
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

	if e.response == nil {
		t.Errorf("The EventListeners response is nil.")
	}

	if e.running {
		t.Errorf("EventListener is already listening although .listen wasn't called.")
	}
}

func TestEventListener_Listen(t *testing.T) {
	setup()
	defer teardown()

	e := Event{"greeting", "Hello, World", "60", time.Now()}

	mux.HandleFunc(eventURL, func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; r.Method != m {
			t.Errorf("Wrong request method %v, expected %v", r.Method, m)
		}

		data, err := json.Marshal(e)

		if err != nil {
			t.Fatalf("Error while encoding event: %v", err)
		}

		fmt.Fprintf(w, ":ok\n\n")
		fmt.Fprintf(w, "event: %v\n", e.Name)
		fmt.Fprintf(w, "data: %v\n\n", string(data[:]))
	})

	eventLister, err := client.NewEventListener("")

	if err != nil {
		t.Fatalf("Error while creating EventLister: %v", err)
	}

	go eventLister.Listen()

	for event := range eventLister.OutputChan {
		if event.Name != e.Name || event.Data != e.Data {
			t.Errorf("Got event %v, expected %v", event, e)
		}
		eventLister.Close()
	}
}

func TestEventListener_ListenDevice(t *testing.T) {
	setup()
	defer teardown()

	d := generateTestDevice("1", "photon", 0)

	e := Event{"greeting", "Hello, World", "60", time.Now()}

	mux.HandleFunc(deviceURL+"/"+d.ID+"/events/"+e.Name, func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; r.Method != m {
			t.Errorf("Wrong request method %v, expected %v", r.Method, m)
		}

		data, err := json.Marshal(e)

		if err != nil {
			t.Fatalf("Error while encoding event: %v", err)
		}

		fmt.Fprintf(w, ":ok\n\n")
		fmt.Fprintf(w, "event: %v\n", e.Name)
		fmt.Fprintf(w, "data: %v\n\n", string(data[:]))
	})

	eventLister, err := d.NewEventListener(e.Name)

	if err != nil {
		t.Fatalf("Error while creating EventLister: %v", err)
	}

	go eventLister.Listen()

	go func() {
		for err := range eventLister.ErrorChan {
			t.Fatalf("Received error from EventListener: %v", err)
		}
	}()

	for event := range eventLister.OutputChan {
		if event.Name != e.Name || event.Data != e.Data {
			t.Errorf("Got event %v, expected %v", event, e)
		}
		eventLister.Close()
	}
}
