package integration

import (
	"http-reverse-proxy/tests/helpers"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadBalancingRoundRobin(t *testing.T) {
	// Initialize logger.
	logger, err := helpers.NewTestLogger()
	assert.NoError(t, err, "Failed to create test logger")

	// Setup two mock backends.
	backendA := helpers.NewMockBackend(200, "Response from Backend A", nil, logger)
	defer backendA.Close()

	backendB := helpers.NewMockBackend(200, "Response from Backend B", nil, logger)
	defer backendB.Close()

	backendURLs := []string{backendA.Server.URL, backendB.Server.URL}

	// Setup proxy server.
	httpServer, teardown := helpers.SetupProxy(t, backendURLs, nil)
	defer teardown()

	// Send multiple requests.
	numRequests := 4
	expectedResponses := []string{
		"Response from Backend A",
		"Response from Backend B",
		"Response from Backend A",
		"Response from Backend B",
	}

	for i := 0; i < numRequests; i++ {
		proxyURL := "http://" + httpServer.Addr + "/loadtest"
		resp, err := helpers.SendRequest("GET", proxyURL, nil)
		assert.NoError(t, err, "Failed to send GET request to proxy")
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err, "Failed to read response body")
		assert.Equal(t, 200, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, expectedResponses[i], string(body), "Unexpected response body")
	}

	// Verify that each backend received the correct number of requests.
	// they each would have received 1 initial health check when the http reverse proxy server is setup
	assert.Len(t, backendA.GetRequests(), 3, "Backend A should receive 2 requests")
	assert.Len(t, backendB.GetRequests(), 3, "Backend B should receive 2 requests")
}
