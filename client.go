package particle

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// A Client manages the communication to the particle cloud.
type Client struct {
	// Http client used to communicate with particle api.
	client *http.Client

	// Base URL for the API requests
	BaseURL *url.URL

	// User agent for the http client.
	UserAgent string

	// Token for authentication.
	Token string
}

// setHeaders sets the authorization and user-agent headers to the given request.
func (c *Client) setHeaders(r *http.Request) {
	r.Header.Add("User-Agent", c.UserAgent)
	r.Header.Add("Authorization", "Bearer "+c.Token)
}

// NewClient returns a new particle cloud api client. If no httpClient was passed,
// than a new one will be created.
func NewClient(httpClient *http.Client, token string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(apiBaseURL)

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
		Token:     token}

	return c
}

// NewRequest creates a new http.Request with the given method to the given endPoint. This function will automatically
// point the request to the clients baseURL, using the clients user agent and token. If a body is passed, than
func (c *Client) NewRequest(method, endPoint string, body io.Reader) (*http.Request, error) {
	// Check that the passed endPoint is valid and concatenate it with the base url.
	path, err := url.Parse(endPoint)

	if err != nil {
		return nil, err
	}

	url := c.BaseURL.ResolveReference(path)

	// Create custom GET request instead of using http.Get so we can headers to it.
	req, err := http.NewRequest(method, url.String(), body)

	if err != nil {
		return nil, err
	}

	c.setHeaders(req)

	return req, nil
}

// Do executes the given http.Request. If the interfaces v is passed, then the function tries to encode the JSON
// response into that interface. The http.Response is passed regardless.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	// Execute the request.
	resp, err := c.client.Do(req)

	if err != nil {
		return resp, err
	}

	// Encode the the JSON response if an interface was passed.
	if v != nil {
		// Be sure to close the body and retrieve any errors.
		defer func() {
			if respErr := resp.Body.Close(); err == nil {
				err = respErr
			}
		}()

		err = CheckResponse(resp)

		if err != nil {
			return nil, err
		}

		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}

// Get executes a GET request using the clients token as well as adding some other headers to it. If the interfaces v is
// passed, then the function tries to encode the JSON response into that interface. The http.Response is passed
// regardless.
func (c *Client) Get(endPoint string, v interface{}) (*http.Response, error) {
	req, err := c.NewRequest("GET", endPoint, nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req, v)

	return resp, err
}

// Post executes a new POST to the given end point with the given form values. If the interfaces v is
// passed, then the function tries to encode the JSON response into that interface. The http.Response is passed
// regardless.
func (c *Client) Post(endPoint string, form url.Values, v interface{}) (*http.Response, error) {
	req, err := c.NewRequest("POST", endPoint, strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(req, v)

	return resp, err
}
