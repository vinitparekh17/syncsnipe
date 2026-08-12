package web

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/server"
	internalsync "github.com/vinitparekh17/syncsnipe/internal/sync"
)

const DefaultPort = "8080"

var port string

func NewWebCmd(frontendDir string, db *database.Queries) (*cobra.Command, error) {
	webCmd := &cobra.Command{
		Use:   "web",
		Short: "run web interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			// signal.NotifyContext cancels ctx on the first SIGINT/SIGTERM, then
			// automatically restores default signal handling so a second Ctrl+C
			// immediately kills the process if graceful shutdown stalls.
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			// Start the file-sync watcher (non-blocking; goroutines exit when ctx is done).
			if db != nil {
				watcher, err := internalsync.NewSyncWatcher(db)
				if err != nil {
					colorlog.Error("failed to create sync watcher: %v", err)
				} else {
					watcher.Start(ctx)
					defer watcher.Close()
				}
			}

			// Build and run the HTTP server; blocks until ctx is cancelled or the
			// server fails. srv.Run does a 5-second graceful shutdown internally.
			srv, err := server.NewServer(port, frontendDir, db)
			if err != nil {
				return fmt.Errorf("failed to create server: %w", err)
			}
			return srv.Run(ctx)
		},
	}

	webCmd.PersistentFlags().StringVarP(&port, "port", "p", DefaultPort, "choose port for web server")
	return webCmd, nil
}
