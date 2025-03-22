package cli

import (
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/cmd/cli/profile"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

func NewCliCmd(q *database.Queries) *cobra.Command {
	cliCmd := &cobra.Command{
		Use:   "cli",
		Short: "run commandline interface",
	}
	cliCmd.AddCommand(profile.NewProfileCmd(q))
	cliCmd.AddCommand(backupCmd)

	return cliCmd
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "snapshot your files for future recovry",
	Run: func(cmd *cobra.Command, args []string) {
		colorlog.Info("backing up......")
	},
}
