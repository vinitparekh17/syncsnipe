package cli

import (
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/cli"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/utils"
)

func NewCliCmd(q *database.Queries) *cobra.Command {
	// root cmd for cli
	cliCmd := &cobra.Command{
		Use:   "cli",
		Short: "run commandline interface",
	}

	// profile management commands
	profileService := core.NewProfile(q)
	profileCmd := NewProfileCmd(profileService)

	// sync-rules management commands
	syncService := core.NewSync(q)
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync your ",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.VerifySuccess(syncService.AddSyncRule(cmd.Context(), args[0], args[1], args[2]), "Sync rule added successfully")
		},
	}

	unsyncCmd := &cobra.Command{
		Use:   "unsync [profile] [sourceDir]",
		Short: "Unsync dir from profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.VerifySuccess(syncService.RemoveSyncRuleByProfile(cmd.Context(), args[0], args[1]), "unsync successfull")
		},
	}

	listSyncCmd := &cobra.Command{
		Use:   "sync-list",
		Short: "List all sync rules with profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			syncRules, err := syncService.ListSyncRules(cmd.Context())
			if err != nil {
				return err
			}
			return cli.DisplayList(syncRules)
		},
	}

	cliCmd.AddCommand(profileCmd)
	cliCmd.AddCommand(syncCmd)
	cliCmd.AddCommand(unsyncCmd)
	cliCmd.AddCommand(listSyncCmd)

	return cliCmd
}
