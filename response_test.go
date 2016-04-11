package particle

import (
	"testing"
)

func TestCheckResponse(t *testing.T) {
	setup()
	defer teardown()


	// Since we define any route we should get an 404 error message
	r, err := client.get("/", nil)

	if err == nil {
		t.Errorf("Recieved no error, but 404 error was expected for request %v", r.Request)
	}

	if r.StatusCode != 404 {
		t.Errorf("Received wrong error code: %v", r.StatusCode)
	}
}
