package particle

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"strconv"
)

const deviceURL = "/v1/devices"

// Device information
type Device struct {
	ID            string
	Name          string
	LastApp       string `json:"last_app"`
	LastIPAddress string `json:"last_ip_address"`
	LastHeard     string `json:"last_heard"`
	ProductID     byte   `json:"product_id"`
	Connected     bool
	Cellular      bool
	Status        string
	LastICCID     string `json:"last_iccid"`
	IMEI          string
	Variables     map[string]string
	Functions     []string
	client        *Client
}

// Devices is an array of the Device type.
type Devices []Device

// FunctionResponse represents the response from the API after calling a device function.
type FunctionResponse struct {
	ID          string
	Name        string
	LastApp     string `json:"last_app"`
	Connected   bool
	ReturnValue int `json:"return_value"`
}

// ListDevices lists the users claimed devices.
func (c *Client) ListDevices() (Devices, error) {
	var devices Devices
	_, err := c.Get(deviceURL, &devices)

	if err != nil {
		return nil, err
	}

	for idx := range devices {
		devices[idx].client = c
	}

	return devices, err
}

// GetDevice gets a single device by it's device
func (c *Client) GetDevice(id string) (Device, error) {
	var device Device
	_, err := c.Get(deviceURL+"/"+id, &device)

	if err != nil {
		return device, err
	}

	device.client = c

	return device, err
}

// variableRaw returns the raw value from a variable as byte buffer for the given device ID and the given variable name.
func (d *Device) variableRaw(name string) (*bytes.Buffer, error) {
	resp, err := d.client.Get(deviceURL+"/"+d.ID+"/"+name+"?format=raw", nil)

	if err != nil {
		return nil, err
	}

	// Be sure to close the body and retrieve any errors.
	defer func() {
		if respErr := resp.Body.Close(); err == nil {
			err = respErr
		}
	}()
	buffer, err := ioutil.ReadAll(resp.Body)

	return bytes.NewBuffer(buffer), err
}

// VariableString returns the string value of the passed devices variable.
func (d *Device) VariableString(name string) (string, error) {
	buffer, err := d.variableRaw(name)
	return buffer.String(), err
}

// VariableInt returns the int value of the passed devices variable.
func (d *Device) VariableInt(name string) (int, error) {
	str, err := d.VariableString(name)

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

// VariableFloat returns the float64 value of the passed devices variable.
func (d *Device) VariableFloat(name string) (float64, error) {
	str, err := d.VariableString(name)

	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(str, 64)
}

// CallFunction calls the passed function name for the given and returns the function value.
func (d *Device) CallFunction(name, argument string) (int, error) {
	form := url.Values{}
	form.Add("arg", argument)
	resp := FunctionResponse{}
	_, err := d.client.Post(deviceURL+"/"+d.ID+"/"+name, form, &resp)

	return resp.ReturnValue, err
}
