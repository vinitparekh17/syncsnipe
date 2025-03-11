package syncsnipe

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/server"
)

var Port string

func NewWebCmd(app *core.App) *cobra.Command {
	return &cobra.Command{
		Use:   "web",
		Short: "run web interface",
		Run: func(cmd *cobra.Command, args []string) {
			server := server.NewServer(app, Port)

			if err := server.Run(); err != nil {
				colorlog.Error("%v", err)
				os.Exit(1)
			}
		},
	}
}
