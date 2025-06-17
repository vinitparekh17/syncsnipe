package web

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/server"
)

const DefaultPort = "8080"

var port string

func NewWebCmd(frontendDir string) (*cobra.Command, error) {
	webCmd := &cobra.Command{
		Use:   "web",
		Short: "run web interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create and start the server
			srv, err := server.NewServer(port, frontendDir)
			if err != nil {
				return fmt.Errorf("failed to create server: %w", err)
			}

			shutdownChan := make(chan struct{})

			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt, syscall.SIGTERM)
				<-c
				close(shutdownChan)
			}()

			if err := srv.Run(shutdownChan); err != nil {
				return fmt.Errorf("server error: %w", err)
			}

			return nil
		},
	}

	webCmd.PersistentFlags().StringVarP(&port, "port", "p", DefaultPort, "choose port for web server")
	return webCmd, nil
}
