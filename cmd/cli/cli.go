package cli

import (
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

// NewCliCmd creates the root CLI command with all its subcommands
func NewCliCmd(q *database.Queries) *cobra.Command {
	cliCmd := &cobra.Command{
		Use:   "cli",
		Short: "Command-line interface for SyncSnipe",
		Long: `SyncSnipe CLI provides commands to manage file synchronization profiles and rules.
Use the available subcommands to create and manage your sync configurations.`,
	}

	// Initialize services
	profileService := core.NewProfile(q)
	syncService := core.NewSync(q)

	// Add subcommands
	cliCmd.AddCommand(
		NewProfileCmd(profileService),
		NewSyncCmd(syncService),
	)

	return cliCmd
}
