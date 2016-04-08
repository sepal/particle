package particle

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

// NewJSONRequest generates a new API request with given request. The urlString should point
// to the API endporint like /v1/devices. An optional body can be passed which is than,
// JSON encoded and send in the request body.
func (c *Client) NewJSONRequest(method, urlString string, body interface{}) (*http.Request, error) {
	path, err := url.Parse(urlString)

	if err != nil {
		return nil, err
	}

	url := c.BaseURL.ResolveReference(path)

	buffer := new(bytes.Buffer)

	if body != nil {
		err := json.NewEncoder(buffer).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url.String(), buffer)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaTypeJSON)
	req.Header.Add("Accept", mediaTypeJSON)
	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Authorization", "Bearer "+c.Token)

	return req, nil
}

// setHeaders sets the authorization and user-agent headers to the given request.
func (c *Client) setHeaders(r *http.Request) {
	r.Header.Add("User-Agent", c.UserAgent)
	r.Header.Add("Authorization", "Bearer "+c.Token)
}

// GET requests executes a GET request using the clients token as well as adding some other headers to it. If v is
// passed it will expect a JSON response from the server and fill the passed interface v with the results. The
// http.Response will be returned either way, as long as there were no errors before the request could be executed.
func (c *Client) Get(endPoint string, v interface{}) (*http.Response, error) {
	// Check that the passed endPoint is valid and concatenate it with the base url.
	path, err := url.Parse(endPoint)

	if err != nil {
		return nil, err
	}

	url := c.BaseURL.ResolveReference(path)

	// Create custom GET request instead of using http.GET so we can headers to it.
	req, err := http.NewRequest("GET", url.String(), nil)

	if err != nil {
		return nil, err
	}

	c.setHeaders(req)

	// If an interface was passed, than we're expecting JSON as a response. Particle unfortunately ignore the
	// request for now :-(
	if v != nil {
		req.Header.Add("Accept", mediaTypeJSON)
	}

	// Execute the request.
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	// Encode the the JSON response if an interface was passed.
	if v != nil {
		// Be sure to close the body and retrieve any errors.
		defer func() {
			if rerr := resp.Body.Close(); err == nil {
				err = rerr
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

// NewFormRequest creates a new Request with form values instead of JSON. The urlString should point
// to the API endporint like /v1/devices.
func (c *Client) NewFormRequest(method, urlString string, form url.Values) (*http.Request, error) {
	path, err := url.Parse(urlString)

	if err != nil {
		return nil, err
	}

	url := c.BaseURL.ResolveReference(path)

	req, err := http.NewRequest(method, url.String(), strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaTypeForm)
	req.Header.Add("Accept", mediaTypeJSON)
	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Authorization", "Bearer "+c.Token)

	return req, nil
}

// Do executes an http.Request and checks for any errors.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	err = CheckResponse(resp)

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(v)

	if err != nil {
		return nil, err
	}

	return resp, err
}

// DoRaw executes a http request with and saves the raw response body into passer buffer element without decoding it.
func (c *Client) DoRaw(req *http.Request, buffer *bytes.Buffer) (*http.Response, error) {
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	buffer.ReadFrom(resp.Body)

	return resp, err
}
