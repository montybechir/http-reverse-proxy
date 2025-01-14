package integration

import (
	"http-reverse-proxy/tests/helpers"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type responseCounter struct {
	healthyA   int
	recoveredB int
}

func TestHealthChecks(t *testing.T) {
	t.Parallel()

	// Initialize logger.
	logger, err := helpers.NewTestLogger()
	assert.NoError(t, err, "Failed to create test logger")

	// Setup two mock backends: one healthy and one initially unhealthy.
	backendA := helpers.NewMockBackend(200, "Healthy Backend A", nil, logger)
	defer backendA.Close()

	backendB := helpers.NewMockBackend(500, "Unhealthy Backend B", nil, logger) // Initially unhealthy
	defer backendB.Close()

	backendURLs := []string{backendA.Server.URL, backendB.Server.URL}

	configOverrides := map[string]interface{}{
		"healthCheckFreq": 2 * time.Second, // Faster health checks for testing
	}
	// Setup proxy server.
	httpServer, teardown := helpers.SetupProxy(t, backendURLs, configOverrides)
	defer teardown()

	// Simulate backendB recovery after some time.
	go func() {
		time.Sleep(1 * time.Second)
		backendB.SetStaticResponse(200, "Recovered Backend B", nil)
	}()

	// Wait to allow health checker to detect changes.
	time.Sleep(20 * time.Second)

	// Send a request to ensure Backend A is responding initially.
	proxyURL := "http://" + httpServer.Addr + "/healthtest"
	resp, err := helpers.SendRequest("GET", proxyURL, nil)
	assert.NoError(t, err, "Failed to send GET request to proxy")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")
	assert.Equal(t, 200, resp.StatusCode, "Expected status code 200")
	assert.Equal(t, "Healthy Backend A", string(body), "Unexpected response from proxy")

	// After recovery, send requests to verify load balancing includes Backend B.
	counter := responseCounter{
		healthyA:   0,
		recoveredB: 0,
	}

	requestCount := 6
	for i := 0; i < requestCount; i++ {
		resp, err := helpers.SendRequest("GET", proxyURL, nil)
		assert.NoError(t, err, "Failed to send GET request to proxy after recovery")
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err, "Failed to read response body after recovery")

		switch string(body) {
		case "Recovered Backend B":
			counter.recoveredB++
		case "Healthy Backend A":
			counter.healthyA++
		default:
			t.Fatalf("Unexpected response: %s", string(body))
		}
	}

	// Log the counts for debugging purposes.
	t.Logf("Response counts after recovery: %+v", counter)

	// Assertions:
	// - Expect at least one response from each backend after recovery.
	// - This accounts for the Round Robin rotation including the recovered backend.
	assert.GreaterOrEqual(t, counter.healthyA, 1, "Expected at least one response from Healthy Backend A after recovery")
	assert.GreaterOrEqual(t, counter.recoveredB, 1, "Expected at least one response from Recovered Backend B after recovery")

	//Further assertions can ensure that the load is roughly balanced.
	// For example, in 6 requests, Backend A and Backend B should each have approximately 3 responses.
	// Allowing a small margin due to timing variations.

	// Calculate expected counts
	expectedA := requestCount / 2
	expectedB := requestCount / 2
	tolerance := 1 // Allowable deviation

	assert.InDelta(t, expectedA, counter.healthyA, float64(tolerance), "Healthy Backend A responses are outside the acceptable range")
	assert.InDelta(t, expectedB, counter.recoveredB, float64(tolerance), "Recovered Backend B responses are outside the acceptable range")

}
