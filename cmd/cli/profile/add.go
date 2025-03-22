package profile

import (
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core/profile"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

func NewAddCmd(q *database.Queries) *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "add a new profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profileName := args[0]
			if err := profile.AddProfile(q, profileName); err != nil {
				colorlog.Fatal("%v", err)
			} else {
				colorlog.Complete("profile %s added successfully", profileName)
			}
		},
	}
}
