package integration

import (
	"http-reverse-proxy/pkg/models"
	"http-reverse-proxy/tests/helpers"
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiting(t *testing.T) {
	t.Parallel()

	// Initialize logger.
	logger, err := helpers.NewTestLogger()
	assert.NoError(t, err, "Failed to create test logger")

	// Setup mock backend.
	backend := helpers.NewMockBackend(200, "Rate Limited Backend", nil, logger)
	defer backend.Close()

	backendURL := backend.Server.URL

	// Override rate limit config for testing (e.g., 1 request per second with a burst of 5).
	testRateLimitConfig := map[string]interface{}{
		"ratelimit": models.RateLimitConfig{
			RequestsPerMinute: 1, // 1 request per minute
			Burst:             5, // Allows up to 5 instantaneous requests
		},
	}

	// Setup proxy server with overridden rate limit.
	httpServer, teardown := helpers.SetupProxy(t, []string{backendURL}, testRateLimitConfig)
	defer teardown()

	// Define the proxy URL endpoint.
	proxyURL := "http://" + httpServer.Addr + "/ratelimittest"

	// Use a WaitGroup to handle concurrent requests if desired. For sequential requests, it's not necessary.
	var wg sync.WaitGroup

	// Define the total number of requests.
	totalRequests := 7

	// Initialize counters for responses.
	counts := map[int]int{
		200: 0,
		429: 0,
	}

	for i := 1; i <= totalRequests; i++ {
		wg.Add(1)
		go func(requestNumber int) {
			defer wg.Done()

			resp, err := helpers.SendRequest("GET", proxyURL, nil)
			if err != nil {
				t.Errorf("Request %d: Failed to send GET request to proxy: %v", requestNumber, err)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Request %d: Failed to read response body: %v", requestNumber, err)
				return
			}

			// Log the response for debugging.
			t.Logf("Request %d: Status Code: %d, Body: %s", requestNumber, resp.StatusCode, string(body))

			// Update counts based on the response status code.
			if resp.StatusCode == 200 {
				counts[200]++
			} else if resp.StatusCode == 429 {
				counts[429]++
			} else {
				t.Errorf("Request %d: Unexpected status code: %d", requestNumber, resp.StatusCode)
			}
		}(i)

		// Introduce a small delay between requests to simulate real-world traffic and to prevent all requests from hitting the burst limit.
		time.Sleep(100 * time.Millisecond) // 0.1 second delay between requests
	}

	// Wait for all requests to complete.
	wg.Wait()

	// Assertions:
	// With RequestsPerMinute = 1 and Burst = 5, we expect:
	// - The first 5 requests to pass (200)
	// - The next 2 requests to be rate limited (429)

	assert.Equal(t, 5, counts[200], "Expected 5 successful requests")
	assert.Equal(t, 2, counts[429], "Expected 2 rate-limited requests")
}
