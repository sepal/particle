package particle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	libraryVersion = "0.0.1"
	apiBaseURL     = "https://api.particle.io"
	userAgent      = "particle/" + libraryVersion
	mediaTypeJSON = "application/json"
	mediaTypeForm = "application/x-www-form-urlencoded"
)

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string
}

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

// CheckResponse checks the API response of an http.Response object.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			return err
		}
	}

	return errorResponse
}
