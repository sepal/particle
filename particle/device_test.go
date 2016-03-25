package particle

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func generateTestDevice(id, name string, productId byte) Device {
	device := Device{
		Id:        id,
		Name:      name,
		ProductID: productId,
		Connected: false,
		Cellular:  productId == 10,
	}

	return device
}

func TestListDevices(t *testing.T) {
	setup()
	defer teardown()

	devices := make(Devices, 2)
	devices[0] = generateTestDevice("1", "core", 0)
	devices[1] = generateTestDevice("1", "electron", 10)

	mux.HandleFunc(deviceUrl, func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		// Return two devices, the first one is a core, the second one an electron

		err := json.NewEncoder(w).Encode(devices)

		if err != nil {
			t.Fatalf("Could not encode devices: %v", err)
		}
	})

	devicesResp, err := client.ListDevices()

	if err != nil {
		t.Fatalf("ListDevices(): %v", err)
	}

	if len(devicesResp) != 2 {
		t.Errorf("Got %v devices, expected 2", len(devicesResp))
	}

	if !reflect.DeepEqual(devicesResp, devices) {
		t.Errorf("Response devices %v don't match with originals: %v", devices, devicesResp)
	}
}
