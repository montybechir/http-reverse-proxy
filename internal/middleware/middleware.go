package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler
type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	// Apply middleware in reverse order
	// Last middleware is executed first
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
