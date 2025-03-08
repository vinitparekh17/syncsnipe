package syncsnipe

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/sync"
)

var rootCmd = &cobra.Command{Use: "syncsnipe"}
var schemaFile = filepath.Join("sql", "schema.sql")


func Execute() {

	db := database.GetDatabase()

	watcher, err := sync.NewSyncWatcher()
	if err != nil {
		colorlog.Fetal("unable to start watcher: %v", err)
	}

	go watcher.Start()

	if err := db.Ping(); err != nil {
		colorlog.Fetal("error while pinging db: %v", err)
	} else {
		if err := db.LoadSchema(schemaFile); err != nil {
			colorlog.Fetal("unable to load schema: %v", err)
		}
		colorlog.Success("Successfully Connected to sqlite")
	}

	dbTx := database.New(db)

	syncSnipeApp := &core.App{
		DBQuery: dbTx,
		Watcher: watcher,
	}

	rootCmd.AddCommand(NewWebCmd(syncSnipeApp))
	rootCmd.AddCommand(NewCliCmd(syncSnipeApp))
	if err := rootCmd.Execute(); err != nil {
		colorlog.Fetal("enable to exec root command: %v", err)
	}
}
