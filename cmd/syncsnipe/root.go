package syncsnipe

import (
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/database"
	s "github.com/vinitparekh17/syncsnipe/internal/sync"
)

const DefaultPort = "8000"

var rootCmd = &cobra.Command{Use: "syncsnipe"}
var schemaFile = filepath.Join("sql", "schema.sql")

func Execute() {
	var wg sync.WaitGroup
	shutdownChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	db := database.GetDatabase()

	if err := db.Ping(); err != nil {
		colorlog.Fetal("error while pinging db: %v", err)
	} else {
		if err := db.LoadSchema(schemaFile); err != nil {
			colorlog.Fetal("unable to load schema: %v", err)
		}
		colorlog.Success("Successfully Connected to sqlite")
	}

	dbTx := database.New(db)

	watcher, err := s.NewSyncWatcher(dbTx)
	if err != nil {
		colorlog.Fetal("unable to start watcher: %v", err)
	}

	app := &core.App{
		DB:           dbTx,
		Watcher:      watcher,
		ShutdownChan: shutdownChan,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		watcher.Start()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-sigChan
		colorlog.Info("shutdown signal received")
		close(shutdownChan)

		watcher.Close()

		if err := db.Close(); err != nil {
			colorlog.Fetal("err closing sqlite3: %v", err)
		}

		colorlog.Success("graceful shutdown completed.")
	}()

	webCmd := NewWebCmd(app)
	webCmd.PersistentFlags().StringVarP(&port, "port", "p", DefaultPort, "choose port for web server")

	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(NewCliCmd(app))
	if err := rootCmd.Execute(); err != nil {
		colorlog.Fetal("unable to exec root command: %v", err)
	}
	wg.Wait()
}
