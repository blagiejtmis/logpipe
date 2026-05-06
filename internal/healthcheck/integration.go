package healthcheck

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// Server wraps an HTTP server that exposes the health endpoint.
type Server struct {
	checker *Checker
	server  *http.Server
	logger  *slog.Logger
}

// NewServer creates a Server that listens on addr and serves /healthz.
func NewServer(addr string, checker *Checker, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	mux := http.NewServeMux()
	mux.Handle("/healthz", checker.Handler())
	return &Server{
		checker: checker,
		server: &http.Server{
			Addr:        addr,
			Handler:     mux,
			ReadTimeout: 5 * time.Second,
		},
		logger: logger,
	}
}

// Run starts the HTTP server and blocks until ctx is cancelled.
func (s *Server) Run(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("healthcheck: listen %s: %w", s.server.Addr, err)
	}
	s.logger.Info("healthcheck server listening", "addr", ln.Addr().String())

	errCh := make(chan error, 1)
	go func() {
		if err := s.server.Serve(ln); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = s.server.Shutdown(shutCtx)
		return nil
	case err := <-errCh:
		return err
	}
}
