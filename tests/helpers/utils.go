// tests/helpers/utils.go

package helpers

import (
	"net/http"
)

// SendRequest sends an HTTP request with optional headers and returns the response.
func SendRequest(method, url string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return client.Do(req)
}
