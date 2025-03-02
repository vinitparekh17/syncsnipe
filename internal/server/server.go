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
)

type SyncSnipe struct{
  server *http.Server

}

func NewServer(mux *http.ServeMux) *SyncSnipe {
  return &SyncSnipe{
    server: &http.Server{
      Addr:         ":8080",
      Handler: mux,
		  ReadTimeout:  5 * time.Second,
    },
  }	
}

func NewMuxRouter() *http.ServeMux {
  mux := http.NewServeMux() 
  return mux
}

func (s *SyncSnipe) Run() error {
  ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
  defer stop()

  errChan := make(chan error, 1)
  go func() {
    colorlog.Info("starting web server on port 8000")
    if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
      errChan <- fmt.Errorf("server error while listening: %v",err)
    }
  }()

  select {
  case <-ctx.Done():
  colorlog.Info("shutting down the server")
  return s.ShutDown(ctx)
  case err := <- errChan:
  return err
  }
}

func (s *SyncSnipe) ShutDown(ctx context.Context) error {
  return s.server.Shutdown(ctx)
  // more clean up task will go here
}
