package profile

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/cli"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core/profile"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

func NewListCmd(q *database.Queries) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list all profiles",
		Run: func(cmd *cobra.Command, args []string) {
			profiles, err := profile.GetProfiles(q)
			if err != nil {
				colorlog.Info("line 20 list.go")
				colorlog.Fatal("%v", err)
			} else {
				if err := cli.DisplayList(profiles); err != nil {
					colorlog.Info("line 23 list.go")
					colorlog.Fatal("%v", err)
				}
				os.Exit(0)
			}
		},
	}
}
