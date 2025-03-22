package profile

import (
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core/profile"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

func NewUpdateCmd(q *database.Queries) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "update your profile",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			oldProfileName := args[0]
			newProfileName := args[1]
			if err := profile.UpdateProfile(q, oldProfileName, newProfileName); err != nil {
				colorlog.Fatal("%v", err)
			} else {
				colorlog.Complete("profile %s updated successfully", newProfileName)
			}
		},
	}
}
