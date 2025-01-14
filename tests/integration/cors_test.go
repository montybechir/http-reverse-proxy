package integration

import (
	"http-reverse-proxy/pkg/models"
	"http-reverse-proxy/tests/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORSPolicies(t *testing.T) {
	// Initialize logger.
	logger, err := helpers.NewTestLogger()
	assert.NoError(t, err, "Failed to create test logger")

	// Setup mock backend.
	backend := helpers.NewMockBackend(200, "CORS Backend", nil, logger)
	defer backend.Close()

	backendURL := backend.Server.URL

	// Override CORS config to allow specific origins.
	testCORSConfig := map[string]interface{}{
		"cors": models.CORSConfig{
			AllowedOrigins: []string{"http://allowed.com", "https://example.com"},
			Debug:          true,
		},
	}

	// Setup proxy server with overridden CORS config.
	httpServer, teardown := helpers.SetupProxy(t, []string{backendURL}, testCORSConfig)
	defer teardown()

	// Define test cases.
	testCases := []struct {
		Name         string
		Origin       string
		ExpectedCode int
	}{
		{
			Name:         "AllowedOrigin1",
			Origin:       "http://allowed.com",
			ExpectedCode: 200,
		},
		{
			Name:         "AllowedOrigin2",
			Origin:       "https://example.com",
			ExpectedCode: 200,
		},
		{
			Name:         "DisallowedOrigin",
			Origin:       "http://disallowed.com",
			ExpectedCode: 403,
		},
		{
			Name:         "NoOrigin",
			Origin:       "",
			ExpectedCode: 403,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			headers := make(map[string]string)
			if tc.Origin != "" {
				headers["Origin"] = tc.Origin
			}
			resp, err := helpers.SendRequest("GET", "http://"+httpServer.Addr+"/corstest", headers)
			assert.NoError(t, err, "Failed to send request")

			assert.Equal(t, tc.ExpectedCode, resp.StatusCode, "Unexpected status code")

			resp.Body.Close()

			// Optionally, verify CORS headers if applicable.
			if tc.ExpectedCode == 200 {
				assert.Equal(t, tc.Origin, resp.Header.Get("Access-Control-Allow-Origin"), "Unexpected CORS header value")
			}
		})
	}
}
