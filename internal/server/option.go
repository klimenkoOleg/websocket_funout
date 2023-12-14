package server

import (
	"net/http"
	"time"
)

type Option func(*options)

type options struct {
	port int
	mux  Mux

	server []func(server *http.Server)
	// mw              []func(handler http.Handler) http.Handler
	shutdownTimeout time.Duration
	logger          Logger
}

// WithReadTimeout sets the maximum duration for reading the entire request, including the body.
func WithReadTimeout(timeout time.Duration) Option {
	return func(o *options) { o.server = append(o.server, func(s *http.Server) { s.ReadTimeout = timeout }) }
}

// WithWriteTimeout sets the maximum duration before timing out writes of the response.
func WithWriteTimeout(timeout time.Duration) Option {
	return func(o *options) { o.server = append(o.server, func(s *http.Server) { s.WriteTimeout = timeout }) }
}

// WithIdleTimeout sets the maximum amount of time to wait for the next request.
func WithIdleTimeout(timeout time.Duration) Option {
	return func(o *options) { o.server = append(o.server, func(s *http.Server) { s.IdleTimeout = timeout }) }
}

// WithShutdownTimeout sets the maximum duration for graceful shutdown (0 is no timeout).
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(o *options) { o.shutdownTimeout = timeout }
}

func WithLogger(logger Logger) Option {
	return func(o *options) { o.logger = logger }
}

func WithPort(port int) Option {
	return func(o *options) { o.port = port }
}
