package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/handler"
	"github.com/vinitparekh17/syncsnipe/internal/stuffbin"
)

type SyncServer struct {
	server *http.Server
	mux    *http.ServeMux
	app    *core.SyncEngine
}

func NewServer(app *core.SyncEngine, port string, frontendDir string) (*SyncServer, error) {
	mux, err := NewMuxRouter(frontendDir)
	if err != nil {
		colorlog.Error("error creating mux router: %v", err)
		return nil, err
	}

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
		mux: mux,
		app: app,
	}, nil
}

func NewMuxRouter(frontendDir string) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	fs, err := stuffbin.LoadFile(frontendDir)
	if err != nil {
		return nil, fmt.Errorf("error loading frontend files: %v", err)
	}

	mux.Handle("/_app/", handler.ServeApp(fs, frontendDir))
	mux.Handle("/assets/", handler.ServeAssets(fs, frontendDir))
	mux.HandleFunc("/", handler.ServeIndexPage(fs, frontendDir))
	return mux, nil
}

func (s *SyncServer) Run(shutDownChan <-chan struct{}) error {
	errChan := make(chan error, 1)

	go func() {
		colorlog.Info("starting web server on port %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server error while listening: %v", err)
		}
	}()

	select {
	case <-shutDownChan:
		colorlog.Info("shutting down web server")
		shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.shutdown(shutDownCtx)
	case err := <-errChan:
		return err
	}
}

func (s *SyncServer) shutdown(shutDownCtx context.Context) error {
	s.app.Watcher.Close()
	return s.server.Shutdown(shutDownCtx)
}
