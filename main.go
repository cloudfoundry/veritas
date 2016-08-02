package main

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/cfdot/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Printf("error: %s\n\n", err.Error())
		commands.RootCmd.Help()
		os.Exit(3)
	}
}
