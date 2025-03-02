package syncsnipe 

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "syncsnipe"}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func init() {
    rootCmd.AddCommand(webCmd)  // From web.go
    rootCmd.AddCommand(cliCmd)  // From cli.go
}
