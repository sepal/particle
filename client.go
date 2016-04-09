package particle

import (
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

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

// Get executes a GET request using the clients token as well as adding some other headers to it. If v is
// passed it will expect a JSON response from the server and fill the passed interface v with the results. The
// http.Response will be returned either way, as long as there were no errors before the request could be executed.
func (c *Client) Get(endPoint string, v interface{}) (*http.Response, error) {
	// Check that the passed endPoint is valid and concatenate it with the base url.
	path, err := url.Parse(endPoint)

	if err != nil {
		return nil, err
	}

	url := c.BaseURL.ResolveReference(path)

	// Create custom GET request instead of using http.Get so we can headers to it.
	req, err := http.NewRequest("GET", url.String(), nil)

	if err != nil {
		return nil, err
	}

	c.setHeaders(req)

	// We're expecting JSON as a response if an interface was passed. Particle unfortunately ignore the request for
	// now :-(
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

// Post executes a new POST to the given end point with the given form values. If v is not null the function will try
// to decode the response as JSON into the give v interface.
func (c *Client) Post(endPoint string, form url.Values, v interface{}) (*http.Response, error) {
	// Check that the passed endPoint is valid and concatenate it with the base url.
	path, err := url.Parse(endPoint)

	if err != nil {
		return nil, err
	}

	url := c.BaseURL.ResolveReference(path)

	// Create custom POST request instead of using http.Post so we can add headers to it.
	req, err := http.NewRequest("POST", url.String(), strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// We're expecting JSON as a response if an interface was passed. Particle unfortunately ignore the request for
	// now :-(
	if v != nil {
		req.Header.Add("Accept", mediaTypeJSON)
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
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
