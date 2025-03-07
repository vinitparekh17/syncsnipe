package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core"
)

type SyncServer struct {
	server *http.Server
	Mux    *http.ServeMux
}

func NewServer(app *core.App) *SyncServer {
	mux := NewMuxRouter()
	return &SyncServer{
		server: &http.Server{
			Addr:        ":8080",
			Handler:     mux,
			ReadTimeout: 5 * time.Second,
		},
		Mux: mux,
	}
}

func NewMuxRouter() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}

func (s *SyncServer) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errChan := make(chan error, 1)
	go func() {
		colorlog.Info("starting web server on port 8000")
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server error while listening: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		colorlog.Info("shutting down the server")
		shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.ShutDown(shutDownCtx)
	case err := <-errChan:
		return err
	}
}

func (s *SyncServer) ShutDown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
	// more clean up task will go here
}
