package main

import (
	"github.com/vinitparekh17/syncsnipe/cmd"
	"github.com/vinitparekh17/syncsnipe/internal/cli"
)

func main() {
	if err := cmd.Execute(); err != nil {
		cli.DisplayError(err)
	}
}
