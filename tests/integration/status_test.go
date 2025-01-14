package integration

import (
	"encoding/json"
	"http-reverse-proxy/tests/helpers"
	"io/ioutil"
	"testing"

	"http-reverse-proxy/internal/proxy"

	"github.com/stretchr/testify/assert"
)

func TestStatusEndpoint(t *testing.T) {
	// Initialize logger.
	logger, err := helpers.NewTestLogger()
	assert.NoError(t, err, "Failed to create test logger")

	// Setup mock backend.
	backend := helpers.NewMockBackend(200, "Status Backend", nil, logger)
	defer backend.Close()

	backendURL := backend.Server.URL

	// Setup proxy server.
	httpServer, teardown := helpers.SetupProxy(t, []string{backendURL}, nil)
	defer teardown()

	// Send request to /status.
	statusURL := "http://" + httpServer.Addr + "/status"
	resp, err := helpers.SendRequest("GET", statusURL, nil)
	assert.NoError(t, err, "Failed to send GET request to /status")
	defer resp.Body.Close()

	// Read and verify response.
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")
	assert.Equal(t, 200, resp.StatusCode, "Expected status code 200")

	// Parse JSON response.
	var statusResp proxy.StatusResponse
	err = json.Unmarshal(body, &statusResp)
	assert.NoError(t, err, "Failed to parse JSON response")
	assert.Equal(t, "running", statusResp.Status, "Unexpected status value")
	assert.NotEmpty(t, statusResp.Uptime, "Uptime should not be empty")
	assert.Equal(t, "1.0.0", statusResp.Version, "Unexpected version value")
}
