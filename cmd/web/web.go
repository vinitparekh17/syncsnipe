package web

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/server"
	s "github.com/vinitparekh17/syncsnipe/internal/sync"
)

var Port string

func NewWebCmd(dbTx *database.Queries) (*cobra.Command, error) {
	var wg sync.WaitGroup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	shutdownChan := make(chan struct{})
	watcher, err := s.NewSyncWatcher(dbTx)
	if err != nil {
		return nil, fmt.Errorf("unable to start watcher: %v", err)
	}

	app := &core.SyncEngine{
		DB:           dbTx,
		Watcher:      watcher,
		ShutdownChan: shutdownChan,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		watcher.Start()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-sigChan
		colorlog.Info("shutdown signal received")
		close(shutdownChan)

		watcher.Close()
		colorlog.Success("graceful shutdown completed.")
	}()

	return &cobra.Command{
		Use:   "web",
		Short: "run web interface",
		Run: func(cmd *cobra.Command, args []string) {
			server := server.NewServer(app, Port)
			if err := server.Run(app.ShutdownChan); err != nil {
				colorlog.Error("%v", err)
				os.Exit(1)
			}
		},
	}, nil
}
