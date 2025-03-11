package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/handler"
	"github.com/vinitparekh17/syncsnipe/internal/stuffbin"
)

type SyncServer struct {
	server *http.Server
	Mux    *http.ServeMux
}

func NewServer(app *core.App, port string) *SyncServer {
	mux := NewMuxRouter()
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		portNum = 8000
	}
	addr := fmt.Sprintf(":%d", portNum)
	return &SyncServer{
		server: &http.Server{
			Addr:        addr,
			Handler:     mux,
			ReadTimeout: 5 * time.Second,
		},
		Mux: mux,
	}
}

func NewMuxRouter() *http.ServeMux {
	mux := http.NewServeMux()
	fs := stuffbin.LoadFile(handler.FrontendDir)
	mux.Handle("/_app/", http.StripPrefix("/_app/", http.FileServer(http.Dir(filepath.Join(handler.FrontendDir, "_app")))))
	mux.HandleFunc("/", handler.ServeIndexPage(fs))
	return mux
}

func (s *SyncServer) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errChan := make(chan error, 1)
	go func() {
		colorlog.Info("starting web server on port %s", s.server.Addr)
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
