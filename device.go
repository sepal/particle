package particle

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
}

// Devices is an array of the Device type.
type Devices []Device

// ListDevices lists the users claimed devices.
func (c *Client) ListDevices() (Devices, error) {
	req, err := c.NewRequest("GET", deviceURL, nil)

	if err != nil {
		return nil, err
	}

	var devices Devices

	_, err = c.Do(req, &devices)

	return devices, err
}

// GetDevice gets a single device by it's device
func (c *Client) GetDevice(id string) (Device, error) {
	req, err := c.NewRequest("GET", deviceURL+"/"+id, nil)

	if err != nil {
		return Device{}, err
	}

	var device Device

	_, err = c.Do(req, &device)

	return device, err
}

func (c *Client) VariableRaw(deviceID string, name string, v interface{}) (error) {
	return nil
}