package proxy

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

// ProxyHandler handles all requests not matched by other routes and proxies them to backends
func (rp *ReverseProxy) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	backendURL, err := rp.LoadBalancer.NextBackend()
	if err != nil {
		rp.Logger.Error("No backend available", zap.Error(err))
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	// Ensure backend URL has scheme
	target := backendURL.Host
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "http://" + target
	}

	// Parse the backend URL
	targetURL, err := url.Parse(target)
	if err != nil {
		rp.Logger.Error("Invalid backend URL",
			zap.String("backend", target),
			zap.Error(err))
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	// Construct full backend URL
	targetURL.Path = r.URL.Path
	targetURL.RawQuery = r.URL.RawQuery

	// Create request to backend
	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		rp.Logger.Error("Failed to create backend request", zap.Error(err))
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	// Copy original headers
	copyHeaders(proxyReq.Header, r.Header)

	// Add proxy headers
	proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
	proxyReq.Header.Set("X-Forwarded-Host", r.Host)
	proxyReq.Header.Set("X-Forwarded-Proto", r.URL.Scheme)
	if r.URL.Scheme == "" {
		proxyReq.Header.Set("X-Forwarded-Proto", "http")
	}

	// Send request to backend
	client := &http.Client{
		Timeout: rp.Config.Server.ReadTimeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
		},
	}

	resp, err := client.Do(proxyReq)
	if err != nil {
		rp.Logger.Error("Backend request failed",
			zap.String("backend", targetURL.String()),
			zap.Error(err))
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	// Stream response body
	written, err := io.Copy(w, resp.Body)
	if err != nil {
		rp.Logger.Error("Failed to copy response body", zap.Error(err))
		return
	}

	rp.Logger.Info("Request proxied successfully",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("backend", targetURL.String()),
		zap.Int64("bytes_written", written),
		zap.Int("status", resp.StatusCode),
	)
}
