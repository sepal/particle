package particle

import (
	"encoding/json"
	"fmt"
	"io"
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
		client:    client,
	}

	return device
}

func TestClient_ListDevices(t *testing.T) {
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

func TestClient_GetDevice(t *testing.T) {
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

func TestDevice_VariableString(t *testing.T) {
	setup()
	defer teardown()

	varName := "message"
	varValue := "My name is particle"

	device := generateTestDevice("1", "core", 0)

	device.Variables = make(map[string]string, 1)
	device.Variables["message"] = "string"

	mux.HandleFunc(deviceURL+"/"+device.ID+"/"+varName, func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}

		io.WriteString(w, varValue)
	})

	response, err := device.VariableString(varName)

	if err != nil {
		t.Fatalf("GetDevice(): %v", err)
	}

	if response != varValue {
		t.Errorf("Variable from response '%v' doesn't match the one generated: '%v'", response, varValue)
	}
}

func TestDevice_VariableInt(t *testing.T) {
	setup()
	defer teardown()

	varName := "anInt"
	varValue := 357

	device := generateTestDevice("1", "core", 0)

	device.Variables = make(map[string]string, 1)
	device.Variables["message"] = "string"

	mux.HandleFunc(deviceURL+"/"+device.ID+"/"+varName, func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}

		io.WriteString(w, fmt.Sprintf("%v", varValue))
	})

	response, err := device.VariableInt(varName)

	if err != nil {
		t.Fatalf("GetDevice(): %v", err)
	}

	if response != varValue {
		t.Errorf("Variable from response '%v' doesn't match the one generated: '%v'", response, varValue)
	}
}

func TestDevice_VariableFloat(t *testing.T) {
	setup()
	defer teardown()

	varName := "anInt"
	varValue := 3.14

	device := generateTestDevice("1", "core", 0)

	device.Variables = make(map[string]string, 1)
	device.Variables["message"] = "string"

	mux.HandleFunc(deviceURL+"/"+device.ID+"/"+varName, func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}

		io.WriteString(w, fmt.Sprintf("%v", varValue))
	})

	response, err := device.VariableFloat(varName)

	if err != nil {
		t.Fatalf("GetDevice(): %v", err)
	}

	if response != float64(varValue) {
		t.Errorf("Variable from response '%v' doesn't match the one generated: '%v'", response, varValue)
	}
}

func TestDevice_CallFunction(t *testing.T) {
	setup()
	defer teardown()

	device := generateTestDevice("1", "photon", 10)
	device.Functions = make([]string, 1)
	device.Functions[0] = "brew"
	funcArg := "coffee"

	mux.HandleFunc(deviceURL+"/"+device.ID+"/"+device.Functions[0], func(w http.ResponseWriter, r *http.Request) {
		if m := "POST"; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}

		err := r.ParseForm()

		if err != nil {
			t.Fatalf("Request body '%v' could not be parsed.", r.Body)
		}

		brew := 1
		if r.PostFormValue("arg") != funcArg {
			t.Errorf("Post form value = %v, expected: %v", r.PostFormValue("arg"), funcArg)
			brew = 0
		}

		resp := FunctionResponse{device.ID, device.Functions[0], "some_app", true, brew}

		err = json.NewEncoder(w).Encode(resp)

		if err != nil {
			t.Fatalf("Could not encode devices: %v", err)
		}
	})

	resp, err := device.CallFunction(device.Functions[0], funcArg)

	if err != nil {
		t.Fatalf("GetDevice(): %v", err)
	}

	if resp != 1 {
		brew := 0

		if funcArg == "coffee" {
			brew = 1
		}

		t.Errorf("Response was '%v', althought '%v' was expected", resp, brew)
	}
}
