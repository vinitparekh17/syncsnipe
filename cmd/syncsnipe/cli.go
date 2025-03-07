package syncsnipe

import (
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/cli"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core"
)
 
func NewCliCmd(app *core.App) *cobra.Command {
	cliCmd := &cobra.Command{
		Use:   "cli",
		Short: "run commandline interface",
	}
	cliCmd.AddCommand(syncCmd)
	cliCmd.AddCommand(backupCmd)

	return cliCmd
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "synchronize files with backup",
	Run: func(cmd *cobra.Command, args []string) {
		cli.SyncDirs(args)
	},
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "snapshot your files for future recovry",
	Run: func(cmd *cobra.Command, args []string) {
		colorlog.Info("backing up......")
	},
}
