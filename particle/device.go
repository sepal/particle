package particle

const deviceUrl = "/v1/devices"

// Device information
type Device struct {
	Id                 string
	Name               string
	LastApp            string `json:"last_app"`
	LastIpAddress      string `json:"last_ip_address"`
	LastHeard          string `json:"last_heard"`
	ProductID          byte   `json:"product_id"`
	Connected          bool
	Cellular           bool
	Status             string
	LastIccid          string `json:"last_iccid"`
	Imei               string
	CurrentBuildTarget string `json:"current_build_target"`
}

// Array of devices
type Devices []Device

// ListDevices() lists the users claimed devices.
func (c *Client) ListDevices() (Devices, error) {
	req, err := c.NewRequest("GET", deviceUrl, nil)

	if err != nil {
		return nil, err
	}

	var devices Devices

	_, err = c.Do(req, &devices)

	return devices, err
}

// Get a single device by it's device
func (c *Client) GetDevice(id string) (Device, error) {
	req, err := c.NewRequest("GET", deviceUrl + "/" + id, nil)

	if err != nil {
		return Device{}, err
	}

	var device Device

	_, err = c.Do(req, &device)

	return device, err
}