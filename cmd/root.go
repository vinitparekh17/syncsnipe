package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/cmd/cli"
	"github.com/vinitparekh17/syncsnipe/cmd/service"
	"github.com/vinitparekh17/syncsnipe/cmd/web"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

var (
	dbFile      = "syncsnipe.db"
	rootCmd     = &cobra.Command{Use: "syncsnipe"}
	schemaFile  = filepath.Join("sql", "schema.sql")
	frontendDir = filepath.Join("frontend", "build")
)

func Execute() error {
	db, err := database.GetDatabase(dbFile)
	if err != nil {
		return fmt.Errorf("unable to get database: %w", err)
	}

	if err := db.LoadSchema(schemaFile); err != nil {
		return fmt.Errorf("unable to load schema: %w", err)
	}
	defer db.Close()

	dbTx := database.New(db)

	webCmd, err := web.NewWebCmd(frontendDir)
	if err != nil {
		return err
	}
	rootCmd.AddCommand(webCmd)

	serviceCmd := service.NewServiceCmd(dbTx)
	rootCmd.AddCommand(serviceCmd)

	cliCmd := cli.NewCliCmd(dbTx)
	rootCmd.AddCommand(cliCmd)

	rootCmd.PersistentFlags().StringVarP(&dbFile, "db", "d", dbFile, "path to the database file")
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("unable to execute root command: %w", err)
	}

	return nil
}
