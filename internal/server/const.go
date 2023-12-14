package server

import (
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"
)

const (
	envPort = "SERVICE_PORT"

	defaultReadTimeout     = 10 * time.Second
	defaultWriteTimeout    = 10 * time.Second
	defaultIdleTimeout     = 60 * time.Second
	defaultShutdownTimeout = 30 * time.Second

	defaultPort = 8080
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGINT, syscall.SIGTERM}

type Mux interface {
	http.Handler
	Handle(string, http.Handler)
}

func envInt(env string, def int) int {
	if p, err := strconv.Atoi(os.Getenv(env)); err == nil {
		return p
	}
	return def
}

func envDuration(env string, def time.Duration) time.Duration {
	if d, err := time.ParseDuration(os.Getenv(env)); err == nil {
		return d
	}
	return def
}
