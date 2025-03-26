package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/cmd/cli"
	"github.com/vinitparekh17/syncsnipe/cmd/web"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

const (
	dbFile = "syncsnipe.db"
)

var rootCmd = &cobra.Command{Use: "syncsnipe"}
var schemaFile = filepath.Join("sql", "schema.sql")

func Execute() error {
	db := database.GetDatabase(dbFile)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging db: %w", err)
	}

	if err := db.LoadSchema(schemaFile); err != nil {
		return fmt.Errorf("unable to load schema: %w", err)
	}

	colorlog.Success("successfully Connected to sqlite")
	defer db.Close()

	dbTx := database.New(db)

	webCmd, err := web.NewWebCmd(dbTx)
	if err != nil {
		return err
	}
	rootCmd.AddCommand(webCmd)

	cliCmd := cli.NewCliCmd(dbTx)
	rootCmd.AddCommand(cliCmd)
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("unable to execute root command: %w", err)
	}

	return nil
}
