package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/cli"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/utils"
)

// NewProfileCmd creates a new profile command with all its subcommands
func NewProfileCmd(profileService core.ProfileService) *cobra.Command {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage sync profiles",
		Long: `Manage sync profiles for your file synchronization tasks.
A profile represents a collection of sync rules that work together to keep your files in sync.`,
	}

	profileCmd.AddCommand(
		newAddProfileCmd(profileService),
		newListProfileCmd(profileService),
		newDeleteProfileCmd(profileService),
		newRenameProfileCmd(profileService),
	)

	return profileCmd
}

// newAddProfileCmd creates the add profile subcommand
func newAddProfileCmd(profileService core.ProfileService) *cobra.Command {
	return &cobra.Command{
		Use:   "add [profile-name]",
		Short: "Add a new sync profile",
		Long: `Add a new sync profile with the specified name.
The profile name must be between 2-50 characters and can only contain letters, numbers, spaces, dashes, and underscores.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			return utils.VerifySuccess(
				profileService.AddProfile(cmd.Context(), profileName),
				"Profile '%s' added successfully",
				profileName,
			)
		},
	}
}

// newListProfileCmd creates the list profiles subcommand
func newListProfileCmd(profileService core.ProfileService) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all sync profiles",
		Long:  "Display a list of all configured sync profiles with their details.",
		RunE: func(cmd *cobra.Command, args []string) error {
			profiles, err := profileService.GetProfiles(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to fetch profiles: %w", err)
			}
			return cli.DisplayList(profiles)
		},
	}
}

// newDeleteProfileCmd creates the delete profile subcommand
func newDeleteProfileCmd(profileService core.ProfileService) *cobra.Command {
	return &cobra.Command{
		Use:   "delete [profile-name]",
		Short: "Delete a sync profile",
		Long: `Delete a sync profile and all its associated sync rules.
This operation cannot be undone.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			return utils.VerifySuccess(
				profileService.DeleteProfile(cmd.Context(), profileName),
				"Profile '%s' deleted successfully",
				profileName,
			)
		},
	}
}

// newRenameProfileCmd creates the rename profile subcommand
func newRenameProfileCmd(profileService core.ProfileService) *cobra.Command {
	return &cobra.Command{
		Use:   "rename [old-name] [new-name]",
		Short: "Rename an existing sync profile",
		Long: `Rename an existing sync profile to a new name.
The new name must be between 2-50 characters and can only contain letters, numbers, spaces, dashes, and underscores.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName, newName := args[0], args[1]
			return utils.VerifySuccess(
				profileService.UpdateProfile(cmd.Context(), oldName, newName),
				"Profile '%s' renamed to '%s'",
				oldName,
				newName,
			)
		},
	}
}
