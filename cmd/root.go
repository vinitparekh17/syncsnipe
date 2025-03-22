package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/cmd/cli"
	"github.com/vinitparekh17/syncsnipe/cmd/web"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

const DefaultPort = "8000"

var rootCmd = &cobra.Command{Use: "syncsnipe"}
var schemaFile = filepath.Join("sql", "schema.sql")

func Execute() {
	db := database.GetDatabase()

	if err := db.Ping(); err != nil {
		colorlog.Fatal("error pinging db: %w", err)
	}

	if err := db.LoadSchema(schemaFile); err != nil {
		colorlog.Fatal("unable to load schema: %w", err)
	}

	colorlog.Success("successfully Connected to sqlite")
	defer db.Close()

	dbTx := database.New(db)

	webCmd := web.NewWebCmd(dbTx)
	webCmd.PersistentFlags().StringVarP(&web.Port, "port", "p", DefaultPort, "choose port for web server")
	rootCmd.AddCommand(webCmd)

	cliCmd := cli.NewCliCmd(dbTx)
	rootCmd.AddCommand(cliCmd)
	rootCmd.AddCommand()
	if err := rootCmd.Execute(); err != nil {
		colorlog.Fatal("unable to exec root command: %v", err)
	}
	// wg.Wait()
}
