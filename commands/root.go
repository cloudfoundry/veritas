package commands

import (
	"code.cloudfoundry.org/lager"
	"github.com/spf13/cobra"
)

var GlobalLogger = lager.NewLogger("cfdot")

var RootCmd = &cobra.Command{
	Use:   "cfdot",
	Short: "Diego operator tooling",
	Long:  "A command-line tool to interact with a Cloud Foundry Diego deployment",
}
