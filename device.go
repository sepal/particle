package particle

import (
	"bytes"
	"strconv"
	"net/url"
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
	Variables map[string]string
	Functions []string
}

// Devices is an array of the Device type.
type Devices []Device

// FunctionResponse represents the response from the API after calling a device function.
type FunctionResponse struct {
	ID string
	Name string
	LastApp string `json:"last_app"`
	Connected bool
	ReturnValue int `json:"return_value"`
}

// ListDevices lists the users claimed devices.
func (c *Client) ListDevices() (Devices, error) {
	req, err := c.NewJSONRequest("GET", deviceURL, nil)

	if err != nil {
		return nil, err
	}

	var devices Devices

	_, err = c.Do(req, &devices)

	return devices, err
}

// GetDevice gets a single device by it's device
func (c *Client) GetDevice(id string) (Device, error) {
	req, err := c.NewJSONRequest("GET", deviceURL+"/"+id, nil)

	if err != nil {
		return Device{}, err
	}

	var device Device

	_, err = c.Do(req, &device)

	return device, err
}

// variableRaw returns the raw value from a variable as byte buffer for the given device ID and the given variable name.
func (c *Client) variableRaw(deviceID, name string) (*bytes.Buffer, error) {
	req, err := c.NewJSONRequest("GET", deviceURL+"/"+deviceID+"/"+name+"?format=raw", nil)

	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	_, err = c.DoRaw(req, buffer)

	return buffer, err
}

// VariableString returns the string value of the passed devices variable.
func (c *Client) VariableString(deviceID, name string) (string, error) {
	buffer, err := c.variableRaw(deviceID, name)
	return buffer.String(), err
}

// VariableInt returns the int value of the passed devices variable.
func (c *Client) VariableInt(deviceID, name string) (int, error) {
	str, err := c.VariableString(deviceID, name)

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

// VariableFloat returns the float64 value of the passed devices variable.
func (c *Client) VariableFloat(deviceID, name string) (float64, error) {
	str, err := c.VariableString(deviceID, name)

	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(str, 64)
}

// CallFunction calls the passed function name for the given and returns the function value.
func (c *Client) CallFunction(deviceID, name, argument string) (int, error) {
	form := url.Values{}
	form.Add("arg", argument)

	req, err := c.NewFormRequest("POST", deviceURL+"/"+deviceID+"/"+name, form)

	if err != nil {
		return 0, err
	}

	resp := FunctionResponse{}
	_, err = c.Do(req, &resp)

	if err != nil {
		return 0, nil
	}

	return resp.ReturnValue, nil
}