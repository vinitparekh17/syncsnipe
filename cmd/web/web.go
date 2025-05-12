package web

import (
	"context"
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

const DefaultPort = "8080"

var port string

func NewWebCmd(dbTx *database.Queries, frontendDir string) (*cobra.Command, error) {
	webCmd := &cobra.Command{
		Use:   "web",
		Short: "run web interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			// Initialize the watcher first to fail fast if there's an error
			watcher, err := s.NewSyncWatcher(dbTx)
			if err != nil {
				return fmt.Errorf("unable to start watcher: %w", err)
			}

			// Setup shutdown signaling
			shutdownChan := make(chan struct{})
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			// Initialize the application
			app := &core.SyncEngine{
				DB:           dbTx,
				Watcher:      watcher,
				ShutdownChan: shutdownChan,
			}

			// Start the server first
			server, err := server.NewServer(app, port, frontendDir)
			if err != nil {
				return fmt.Errorf("failed to create server: %w", err)
			}

			var wg sync.WaitGroup

			// Start the watcher
			wg.Add(1)
			go func() {
				defer wg.Done()
				watcher.Start(ctx)
			}()

			// Handle shutdown signals
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case <-sigChan:
					colorlog.Info("shutdown signal received")
					close(shutdownChan)
				case <-ctx.Done():
					colorlog.Info("context cancelled")
				}
			}()

			// Run the server - this will block until shutdown
			if err := server.Run(app.ShutdownChan); err != nil {
				colorlog.Error("server error: %v", err)
				return fmt.Errorf("server error: %w", err)
			}

			// Wait for all goroutines to finish
			wg.Wait()
			return nil
		},
	}

	webCmd.PersistentFlags().StringVarP(&port, "port", "p", DefaultPort, "choose port for web server")
	return webCmd, nil
}
