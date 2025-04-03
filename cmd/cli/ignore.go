package cli

import (
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/cli"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/utils"
)

// NewIgnorePatternCmd creates a new ignore pattern command with all its subcommands
func NewIgnorePatternCmd(ignoreService core.IgnoreService) *cobra.Command {
	ignorePatternCmd := &cobra.Command{
		Use:   "ignore",
		Short: "Manage ignore patterns",
	}

	ignorePatternCmd.AddCommand(
		newAddIgnorePatternCmd(ignoreService),
		newListIgnorePatternCmd(ignoreService),
		newDeleteIgnorePatternCmd(ignoreService),
	)

	return ignorePatternCmd
}

// newAddIgnorePatternCmd creates a new add ignore pattern command
func newAddIgnorePatternCmd(ignoreService core.IgnoreService) *cobra.Command {
	addIgnorePatternCmd := &cobra.Command{
		Use:   "add [profile-name] [pattern]",
		Short: "Add an ignore pattern to a profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.VerifySuccess(ignoreService.AddIgnore(cmd.Context(), args[0], args[1]), "ignore pattern added successfully")
			return nil
		},
	}

	return addIgnorePatternCmd
}

// newListIgnorePatternCmd creates a new list ignore pattern command
func newListIgnorePatternCmd(ignoreService core.IgnoreService) *cobra.Command {
	listIgnorePatternCmd := &cobra.Command{
		Use:   "list [profile-name]",
		Short: "List all ignore patterns for a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ignorePatterns, err := ignoreService.ListIgnorePatterns(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			return cli.DisplayList(ignorePatterns)
		},
	}

	return listIgnorePatternCmd
}

// newDeleteIgnorePatternCmd creates a new delete ignore pattern command
func newDeleteIgnorePatternCmd(ignoreService core.IgnoreService) *cobra.Command {
	deleteIgnorePatternCmd := &cobra.Command{
		Use:   "remove [profile-name] [pattern]",
		Short: "Remove an ignore pattern from a profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.VerifySuccess(ignoreService.DeleteIgnorePattern(cmd.Context(), args[0], args[1]), "ignore pattern deleted successfully")
			return nil
		},
	}

	return deleteIgnorePatternCmd
}
