package profile

import (
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

func NewProfileCmd(q *database.Queries) *cobra.Command {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "manage profiles",
	}

	profileCmd.AddCommand(NewAddCmd(q))
	profileCmd.AddCommand(NewListCmd(q))
	profileCmd.AddCommand(NewDeleteCmd(q))
	profileCmd.AddCommand(NewUpdateCmd(q))

	return profileCmd
}
