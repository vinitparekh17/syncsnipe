package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime"
	"strconv"
	"time"

	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/handler"
	"github.com/vinitparekh17/syncsnipe/internal/stuffbin"
)

type SyncServer struct {
	server *http.Server
	mux    *http.ServeMux
}

func NewServer(port string, frontendDir string, db *database.Queries) (*SyncServer, error) {
	mux, err := NewMuxRouter(frontendDir, db)
	if err != nil {
		colorlog.Error("error creating mux router: %v", err)
		return nil, err
	}

	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		portNum = 8080
	}
	addr := fmt.Sprintf(":%d", portNum)
	return &SyncServer{
		server: &http.Server{
			Addr:        addr,
			Handler:     mux,
			ReadTimeout: 10 * time.Second,
		},
		mux: mux,
	}, nil
}

func NewMuxRouter(frontendDir string, db *database.Queries) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	// Enable block and mutex profiling for concurrency analysis
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	// Register pprof handlers
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	if db != nil {
		apiHandler := handler.NewAPIHandler(db)
		mux.HandleFunc("/api/dashboard/stats", apiHandler.GetDashboardStats)
		mux.HandleFunc("/api/profiles", apiHandler.HandleProfiles)
		mux.HandleFunc("/api/sync-rules", apiHandler.HandleSyncRules)
		mux.HandleFunc("/api/activities", apiHandler.GetRecentActivities)
	}

	fs, err := stuffbin.LoadFile(frontendDir)
	if err != nil {
		return nil, fmt.Errorf("error loading frontend files: %v", err)
	}

	mux.Handle("/_app/", handler.ServeApp(fs, frontendDir))
	mux.HandleFunc("/", handler.ServeIndexPage(fs, frontendDir))
	return mux, nil
}

// Run starts the HTTP server and blocks until ctx is cancelled or the server fails.
func (s *SyncServer) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		colorlog.Info("starting web server on port %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server error while listening: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		colorlog.Info("shutting down web server")
		shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.shutdown(shutDownCtx)
	case err := <-errChan:
		return err
	}
}

func (s *SyncServer) shutdown(shutDownCtx context.Context) error {
	return s.server.Shutdown(shutDownCtx)
}
