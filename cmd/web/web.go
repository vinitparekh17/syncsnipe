package web

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/server"
)

const DefaultPort = "8080"

var port string

func NewWebCmd(frontendDir string) (*cobra.Command, error) {
	webCmd := &cobra.Command{
		Use:   "web",
		Short: "run web interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			shutdownChan := make(chan struct{})

			// Start the server first
			server, err := server.NewServer(port, frontendDir)
			if err != nil {
				return fmt.Errorf("failed to create server: %w", err)
			}

			// Run the server - this will block until shutdown
			if err := server.Run(shutdownChan); err != nil {
				colorlog.Error("server error: %v", err)
				return fmt.Errorf("server error: %w", err)
			}

			return nil
		},
	}

	webCmd.PersistentFlags().StringVarP(&port, "port", "p", DefaultPort, "choose port for web server")
	return webCmd, nil
}
