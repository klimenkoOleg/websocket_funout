package server

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type HTTPListener struct {
	options
}

func New(opts ...Option) *HTTPListener {
	o := options{
		shutdownTimeout: defaultShutdownTimeout,
		port:            envInt(envPort, defaultPort),
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &HTTPListener{
		options: o,
	}
}

func (s *HTTPListener) Listen(ctx context.Context, handler http.Handler) error {

	server := &http.Server{
		Addr:           ":" + strconv.Itoa(s.port),
		Handler:        handler,
		ReadTimeout:    defaultReadTimeout,
		WriteTimeout:   defaultWriteTimeout,
		IdleTimeout:    defaultIdleTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ErrorLog:       errorLog(s.logger),
	}

	for _, opt := range s.server {
		opt(server)
	}

	chErrors := make(chan error)
	chSignals := make(chan os.Signal, 2)
	signal.Notify(chSignals, shutdownSignals...)

	go listen(chErrors, server)

	var err error
	select {
	case err = <-chErrors:
		_ = shutdown(server, s.shutdownTimeout)
	case <-chSignals:
		signal.Stop(chSignals)
		err = shutdown(server, s.shutdownTimeout)
	case <-ctx.Done():
		err = ctx.Err()
		if e := shutdown(server, s.shutdownTimeout); e != nil {
			err = e
		}
	}

	close(chErrors)
	close(chSignals)

	return err
}

func listen(ch chan error, server *http.Server) {

	err := server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		ch <- err
	}
}

func shutdown(server *http.Server, timeout time.Duration) error {
	var cancel func()
	ctx := context.Background()

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	server.SetKeepAlivesEnabled(false)
	return server.Shutdown(ctx)
}

func errorLog(logger Logger) *log.Logger {
	if logger == nil {
		return log.New(os.Stderr, "", log.LstdFlags)
	}

	return log.New(writerFunc(func(p []byte) (n int, err error) {
		l := len(p)

		if bytes.HasPrefix(p, []byte("http: panic serving ")) {
			// Skip logging of panic, handled by server itself. This kind of errors would be logged
			// by middlewares
			return l, nil
		}

		p = bytes.TrimRight(p, "\n")

		logger.Error(context.Background(), string(p))

		return l, nil
	}), "", 0)
}
