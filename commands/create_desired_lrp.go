package commands

import (
	"encoding/json"
	"fmt"
	"io"

	"code.cloudfoundry.org/bbs"
	"code.cloudfoundry.org/bbs/models"
	"github.com/spf13/cobra"
)

var createDesiredLRPCmd = &cobra.Command{
	Use:   "create-desired-lrp",
	Short: "Create a desired LRP",
	Long:  "Create a desired LRP from the given specs",
	RunE:  createDesiredLRP,
}

func init() {
	AddBBSFlags(createDesiredLRPCmd)
	RootCmd.AddCommand(createDesiredLRPCmd)
}

func createDesiredLRP(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected one argument, found %d", len(args))
	}

	bbsClient, err := newBBSClient(cmd)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	err = CreateDesiredLRP(cmd, cmd.OutOrStdout(), cmd.OutOrStderr(), bbsClient, args[0])
	if err != nil {
		return err
	}

	return nil
}

func CreateDesiredLRP(cmd *cobra.Command, stdout, stderr io.Writer, bbsClient bbs.Client, spec string) error {
	logger := globalLogger.Session("desiredLRPs")

	var desiredLRP *models.DesiredLRP
	err := json.Unmarshal([]byte(spec), &desiredLRP)
	if err != nil {
		return NewCFDotValidationError(cmd, err)
	}

	err = bbsClient.DesireLRP(logger, desiredLRP)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	return nil
}
