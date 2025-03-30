package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/cli"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/utils"
)

// NewSyncCmd creates a new sync command with all its subcommands
func NewSyncCmd(syncService core.SyncService) *cobra.Command {
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Manage file synchronization rules",
		Long: `Manage file synchronization rules for your profiles.
A sync rule defines how files from a source directory should be synchronized to a target directory.`,
	}

	syncCmd.AddCommand(
		newAddSyncCmd(syncService),
		newUnsyncCmd(syncService),
		newListSyncCmd(syncService),
		newStatusSyncCmd(syncService),
	)

	return syncCmd
}

// newAddSyncCmd creates the add sync rule subcommand
func newAddSyncCmd(syncService core.SyncService) *cobra.Command {
	return &cobra.Command{
		Use:   "add [profile] [source-dir] [target-dir]",
		Short: "Add a new sync rule",
		Long: `Add a new sync rule to a profile.
This will start monitoring the source directory and sync its contents to the target directory.
The source and target directories must be different.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, sourceDir, targetDir := args[0], args[1], args[2]
			return utils.VerifySuccess(
				syncService.AddSyncRule(cmd.Context(), profile, sourceDir, targetDir),
				"Sync rule added successfully for profile '%s'",
				profile,
			)
		},
	}
}

// newUnsyncCmd creates the unsync subcommand
func newUnsyncCmd(syncService core.SyncService) *cobra.Command {
	return &cobra.Command{
		Use:   "remove [profile] [source-dir]",
		Short: "Remove a sync rule",
		Long: `Remove a sync rule from a profile.
This will stop monitoring the specified source directory for this profile.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, sourceDir := args[0], args[1]
			return utils.VerifySuccess(
				syncService.RemoveSyncRuleByProfile(cmd.Context(), profile, sourceDir),
				"Sync rule removed successfully from profile '%s'",
				profile,
			)
		},
	}
}

// newListSyncCmd creates the list sync rules subcommand
func newListSyncCmd(syncService core.SyncService) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all sync rules",
		Long:  "Display a list of all configured sync rules grouped by profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			syncRules, err := syncService.ListSyncRules(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to fetch sync rules: %w", err)
			}
			return cli.DisplayList(syncRules)
		},
	}
}

func newStatusSyncCmd(syncService core.SyncService) *cobra.Command {
	return &cobra.Command{
		Use:   "status [profile]",
		Short: "Check the status of a sync rule",
		Long:  "Check the status of a sync rule for a specific profile.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			syncStatus, err := syncService.GetSyncStatusByProfileName(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("failed to get sync status: %w", err)
			}
			return cli.DisplayStruct(syncStatus)
		},
	}
}
