package cli

import (
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/cli"
	"github.com/vinitparekh17/syncsnipe/internal/core/profile"
)

func NewProfileCmd(profileService profile.ProfileService) *cobra.Command {
	// Helper function to create commands with consistent error handling
	createCmd := func(use, short string, args int, run func(cmd *cobra.Command, args []string) error) *cobra.Command {
		return &cobra.Command{
			Use:   use,
			Short: short,
			Args:  cobra.ExactArgs(args),
			RunE:  run,
		}
	}

	addProfileCmd := createCmd(
		"add [profile-name]",
		"Add a new profile",
		1,
		func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			return profileService.AddProfile(cmd.Context(), profileName)
		},
	)

	listProfileCmd := createCmd(
		"list",
		"List all profiles",
		0,
		func(cmd *cobra.Command, args []string) error {
			profiles, err := profileService.GetProfiles(cmd.Context())
			if err != nil {
				return err
			}
			return cli.DisplayList(profiles)
		},
	)

	deleteProfileCmd := createCmd(
		"delete [profile-name]",
		"Delete a profile",
		1,
		func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			return profileService.DeleteProfile(cmd.Context(), profileName)
		},
	)

	editProfileCmd := createCmd(
		"rename [old-name] [new-name]",
		"Rename an existing profile",
		2,
		func(cmd *cobra.Command, args []string) error {
			oldName := args[0]
			newName := args[1]
			return profileService.UpdateProfile(cmd.Context(), oldName, newName)
		},
	)

	// Parent command to group all profile commands
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage profiles",
	}
	profileCmd.AddCommand(addProfileCmd, listProfileCmd, deleteProfileCmd, editProfileCmd)

	return profileCmd
}
