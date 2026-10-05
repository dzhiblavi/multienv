package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"multienv/server/internal/dispatcher"
)

type Server struct {
	mux    *http.ServeMux
	server *http.Server
	config Config
}

func New(config Config) Server {
	mux := http.NewServeMux()

	return Server{
		mux,
		&http.Server{
			Addr:    fmt.Sprintf("%s:%d", config.Address, config.Port),
			Handler: mux,
		},
		config,
	}
}

func (s Server) Serve(ctx context.Context) error {
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- s.server.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		_ = s.server.Close()
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	err := <-serveErr
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (s Server) Install(name string, disp dispatcher.Dispatcher) error {
	s.mux.HandleFunc("POST /"+name, makeHandler(disp))
	return nil
}
