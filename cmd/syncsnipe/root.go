package syncsnipe

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/sync"
)

const DefaultPort = "8000"

var rootCmd = &cobra.Command{Use: "syncsnipe"}
var schemaFile = filepath.Join("sql", "schema.sql")

var app *core.App
var webCmd = NewWebCmd(app)

func init() {
	webCmd.PersistentFlags().StringVarP(&Port, "port", "p", DefaultPort, "choose port for web server")
}

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

	app = &core.App{
		DBQuery: dbTx,
		Watcher: watcher,
	}
	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(NewCliCmd(app))
	if err := rootCmd.Execute(); err != nil {
		colorlog.Fetal("enable to exec root command: %v", err)
	}
}
