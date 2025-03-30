package main

import "github.com/vinitparekh17/syncsnipe/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
