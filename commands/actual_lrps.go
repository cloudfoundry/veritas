package commands

import (
	"encoding/json"
	"io"

	"code.cloudfoundry.org/bbs"
	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/cfdot/commands/helpers"
	"github.com/spf13/cobra"
)

// flags
var (
	actualLRPsDomainFlag, actualLRPsCellIdFlag string
)

var actualLRPsCmd = &cobra.Command{
	Use:   "actual-lrps",
	Short: "List actual LRPs",
	Long:  "List actual LRPs from the BBS",
	RunE:  actualLRPs,
}

func init() {
	AddBBSAndTimeoutFlags(actualLRPsCmd)

	actualLRPsCmd.Flags().StringVarP(&actualLRPsDomainFlag, "domain", "d", "", "retrieve only actual lrps for the given domain")
	actualLRPsCmd.Flags().StringVarP(&actualLRPsCellIdFlag, "cell-id", "c", "", "retrieve only actual lrps for the given cell id")

	RootCmd.AddCommand(actualLRPsCmd)
}

func actualLRPs(cmd *cobra.Command, args []string) error {
	err := ValidateActualLRPsArguments(args)
	if err != nil {
		return NewCFDotValidationError(cmd, err)
	}

	bbsClient, err := helpers.NewBBSClient(cmd, Config)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	err = ActualLRPs(
		cmd.OutOrStdout(),
		cmd.OutOrStderr(),
		bbsClient,
		actualLRPsDomainFlag,
		actualLRPsCellIdFlag,
	)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	return nil
}

func ValidateActualLRPsArguments(args []string) error {
	if len(args) > 0 {
		return errExtraArguments
	}
	return nil
}

func ActualLRPs(stdout, stderr io.Writer, bbsClient bbs.Client, domain, cellID string) error {
	logger := globalLogger.Session("actual-lrps")

	encoder := json.NewEncoder(stdout)

	actualLRPFilter := models.ActualLRPFilter{
		CellID: cellID,
		Domain: domain,
	}

	actualLRPs, err := bbsClient.ActualLRPs(logger, actualLRPFilter)
	if err != nil {
		return err
	}

	for _, actualLRP := range actualLRPs {
		err = encoder.Encode(actualLRP)
		if err != nil {
			logger.Error("failed-to-marshal", err)
		}
	}

	return nil
}
