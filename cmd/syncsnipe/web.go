package syncsnipe

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/handler"
	"github.com/vinitparekh17/syncsnipe/internal/server"
	"github.com/vinitparekh17/syncsnipe/internal/stuffbin"
)

var (
	staticFile = filepath.Join("frontend", "build", "index.html")
)

func NewWebCmd(app *core.App) *cobra.Command {
	return &cobra.Command{
		Use:   "web",
		Short: "run web interface",
		Run: func(cmd *cobra.Command, args []string) {
			server := server.NewServer(app)
			fs := stuffbin.LoadFile(staticFile)
			server.Mux.Handle("/", handler.HandleStaticFiles(staticFile, fs))
      if err := server.Run(); err != nil {
				colorlog.Error("%v", err)
				os.Exit(1)
			}
		},
	}
}
