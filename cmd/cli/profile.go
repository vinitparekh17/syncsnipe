package cli

import (
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/cli"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/utils"
)

func NewProfileCmd(profileService core.ProfileService) *cobra.Command {

	addProfileCmd := &cobra.Command{
		Use:   "add [profile-name]",
		Short: "Add a new profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			return utils.VerifySuccess(profileService.AddProfile(cmd.Context(), profileName), "%s profile added successfully", profileName)
		},
	}

	listProfileCmd := &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			profiles, err := profileService.GetProfiles(cmd.Context())
			if err != nil {
				return err
			}
			return cli.DisplayList(profiles)
		},
	}

	deleteProfileCmd := &cobra.Command{
		Use:   "delete [profile-name]",
		Short: "Delete a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			return utils.VerifySuccess(profileService.DeleteProfile(cmd.Context(), profileName), "%s profile deleted successfully", profileName)
		},
	}

	editProfileCmd := &cobra.Command{
		Use:   "rename [old-name] [new-name]",
		Short: "Rename an existing profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName := args[0]
			newName := args[1]
			return utils.VerifySuccess(profileService.UpdateProfile(cmd.Context(), oldName, newName), "%s profile renamed to %s", oldName, newName)
		},
	}

	// Parent command to group all profile commands
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage profiles",
	}
	profileCmd.AddCommand(addProfileCmd, listProfileCmd, deleteProfileCmd, editProfileCmd)

	return profileCmd
}
