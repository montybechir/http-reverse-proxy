package integration

import (
	"http-reverse-proxy/tests/helpers"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyFunctionality(t *testing.T) {
	// Initialize logger (assuming you have a helper for this; otherwise modify accordingly)
	logger, err := helpers.NewTestLogger()
	assert.NoError(t, err, "Failed to create test logger")

	// Setup mock backend.
	backend := helpers.NewMockBackend(200, "Hello from Backend A", nil, logger)
	defer backend.Close()

	backendURL := backend.Server.URL

	// Setup proxy server.
	httpServer, teardown := helpers.SetupProxy(t, []string{backendURL}, nil)
	defer teardown()

	// Send a request to the proxy.
	proxyURL := "http://" + httpServer.Addr + "/testpath"
	resp, err := helpers.SendRequest("GET", proxyURL, nil)
	assert.NoError(t, err, "Failed to send GET request to proxy")
	defer resp.Body.Close()

	// Read response body.
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")
	assert.Equal(t, 200, resp.StatusCode, "Expected status code 200")
	assert.Equal(t, "Hello from Backend A", string(body), "Unexpected response body")

	// Ensure the backend received the request.
	requests := backend.GetRequests()

	// Since our proxy server initializes a test to backends when it starts, we expect 2 requests
	assert.Len(t, requests, 2, "Expected one request to backend")

	// the first request will be the one sent by the proxy itself, the second request is the /testpath req above
	assert.Equal(t, "/testpath", requests[1].URL.Path, "Unexpected request path")
}
