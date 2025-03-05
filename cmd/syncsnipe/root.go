package syncsnipe

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/sync"
)

var rootCmd = &cobra.Command{Use: "syncsnipe"}

func Execute() {

	db := database.GetDatabase()
	watcher, err := sync.NewSyncWatcher()
	if err != nil {
		colorlog.Error("unable to start watcher: %v", err)
		os.Exit(1)
	}

  go watcher.Start()
  if err := db.Ping(); err != nil {
    colorlog.Error("error while pinging db: %v", err)
    os.Exit(1)
  } else {
    colorlog.Success("Successfully Connected to sqlite")
  }
	syncSnipeApp := &core.App{
		DB:      db,
		Watcher: watcher,
  }

	rootCmd.AddCommand(NewWebCmd(syncSnipeApp))
	rootCmd.AddCommand(NewCliCmd(syncSnipeApp))
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
