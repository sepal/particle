package particle

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string
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
