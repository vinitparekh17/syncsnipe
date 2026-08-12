package cmd

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/cmd/cli"
	"github.com/vinitparekh17/syncsnipe/cmd/service"
	"github.com/vinitparekh17/syncsnipe/cmd/web"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

var (
	dbFile      = "syncsnipe.db"
	pprofAddr   = ""
	rootCmd     = &cobra.Command{Use: "syncsnipe"}
	schemaFile  = filepath.Join("sql", "schema.sql")
	frontendDir = filepath.Join("frontend", "build")
)

func Execute() error {
	// Windows Service Control Manager launches services from a system working
	// directory. Keep all bundled assets and the default database next to the
	// executable so foreground and background runs use the same configuration.
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("unable to resolve executable path: %w", err)
	}
	if err := os.Chdir(filepath.Dir(execPath)); err != nil {
		return fmt.Errorf("unable to set application working directory: %w", err)
	}

	db, err := database.GetDatabase(dbFile)
	if err != nil {
		return fmt.Errorf("unable to get database: %w", err)
	}

	if err := db.LoadSchema(schemaFile); err != nil {
		return fmt.Errorf("unable to load schema: %w", err)
	}
	defer db.Close()

	dbTx := database.New(db)

	webCmd, err := web.NewWebCmd(frontendDir, dbTx)
	if err != nil {
		return err
	}
	rootCmd.AddCommand(webCmd)

	serviceCmd := service.NewServiceCmd(dbTx)
	rootCmd.AddCommand(serviceCmd)

	cliCmd := cli.NewCliCmd(dbTx)
	rootCmd.AddCommand(cliCmd)

	rootCmd.PersistentFlags().StringVarP(&dbFile, "db", "d", dbFile, "path to the database file")
	rootCmd.PersistentFlags().StringVar(&pprofAddr, "pprof", "", "enable dedicated pprof profiling server on address (e.g. :6060)")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if pprofAddr != "" {
			runtime.SetBlockProfileRate(1)
			runtime.SetMutexProfileFraction(1)
			go func() {
				colorlog.Info("starting dedicated pprof profiling server on %s", pprofAddr)
				if err := http.ListenAndServe(pprofAddr, nil); err != nil && err != http.ErrServerClosed {
					colorlog.Error("pprof server error: %v", err)
				}
			}()
		}
	}

	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("unable to execute root command: %w", err)
	}

	return nil
}
