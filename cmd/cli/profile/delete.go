package profile

import (
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core/profile"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

func NewDeleteCmd(q *database.Queries) *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "delete a profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profileName := args[0]
			if err := profile.DeleteProfile(q, profileName); err != nil {
				colorlog.Fatal("%v", err)
			} else {
				colorlog.Complete("profile %s has been deleted successfully", profileName)
			}
		},
	}
}
