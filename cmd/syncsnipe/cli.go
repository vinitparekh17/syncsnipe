package syncsnipe

import (
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
)

var cliCmd = &cobra.Command{
  Use: "cli",
  Short: "run commandline interface",
  Run: func(cmd *cobra.Command, args []string) {
    colorlog.Success("running SyncSnipe in cli mode")
  },
}
