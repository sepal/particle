package particle

const deviceUrl = "/v1/devices"

type Device struct {
	Id             string
	Name           string
	LastApp        string
	LastIp_address string
	LastHeard      string
	ProductID      byte
	LastIccid      string
	Imei           string
}

func (c *Client) ListDevices() ([]Device, error) {
	req, err := c.NewRequest("GET", deviceUrl, nil)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
