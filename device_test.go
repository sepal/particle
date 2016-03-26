package particle

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

// generateTestDevice generates a device for testing.
func generateTestDevice(id, name string, productID byte) Device {
	device := Device{
		ID:        id,
		Name:      name,
		ProductID: productID,
		Connected: false,
		Cellular:  productID == 10,
	}

	return device
}

func TestListDevices(t *testing.T) {
	setup()
	defer teardown()

	// Generate two devices, the first one is a core and the second one an electron.
	devices := make(Devices, 2)
	devices[0] = generateTestDevice("1", "core", 0)
	devices[1] = generateTestDevice("1", "electron", 10)

	mux.HandleFunc(deviceURL, func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}

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
		t.Errorf("Response devices %v don't match with originals: %v", devicesResp, devices)
	}
}

func TestGetDevice(t *testing.T) {
	setup()
	defer teardown()

	device := generateTestDevice("1", "core", 0)

	mux.HandleFunc(deviceURL+"/"+device.ID, func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}

		err := json.NewEncoder(w).Encode(device)

		if err != nil {
			t.Fatalf("Could not encode device: %v", err)
		}
	})

	deviceResp, err := client.GetDevice(device.ID)

	if err != nil {
		t.Fatalf("GetDevice(): %v", err)
	}

	if !reflect.DeepEqual(deviceResp, device) {
		t.Errorf("Response device %v doesn't match orignal: %v", deviceResp, device)
	}
}
