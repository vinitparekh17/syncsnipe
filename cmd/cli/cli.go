package cli

import (
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/core/profile"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

func NewCliCmd(q *database.Queries) *cobra.Command {
	cliCmd := &cobra.Command{
		Use:   "cli",
		Short: "run commandline interface",
	}

	profileService := profile.NewProfile(q)
	cliCmd.AddCommand(NewProfileCmd(profileService))

	return cliCmd
}
