package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"code.cloudfoundry.org/bbs"

	"code.cloudfoundry.org/cfdot/commands/helpers"
	"github.com/spf13/cobra"
)

// flags
var (
	actualLRPsGuidIndexFlag string
)

var actualLRPsByProcessGuidCmd = &cobra.Command{
	Use:   "actual-lrps-for-guid PROCESS_GUID",
	Short: "List actual LRPs for a process guid",
	Long:  fmt.Sprintf("List actual LRPs from the BBS for a given process guid. Process guids can be obtained by running %s actual-lrps", os.Args[0]),
	RunE:  actualLRPsByProcessGuid,
}

func init() {
	AddBBSAndTimeoutFlags(actualLRPsByProcessGuidCmd)

	// String flag because logic for optional int flag is not clear
	actualLRPsByProcessGuidCmd.Flags().StringVarP(&actualLRPsGuidIndexFlag, "index", "i", "", "retrieve actual lrp for the given index")

	RootCmd.AddCommand(actualLRPsByProcessGuidCmd)
}

func actualLRPsByProcessGuid(cmd *cobra.Command, args []string) error {
	processGuid, index, err := ValidateActualLRPsForGuidArgs(args, actualLRPsGuidIndexFlag)
	if err != nil {
		return NewCFDotValidationError(cmd, err)
	}

	bbsClient, err := helpers.NewBBSClient(cmd, Config)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	err = ActualLRPsForGuid(cmd.OutOrStdout(), cmd.OutOrStderr(), bbsClient, processGuid, index)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	return nil
}

func ValidateActualLRPsForGuidArgs(args []string, indexFlag string) (string, int, error) {
	if len(args) < 1 {
		return "", 0, errMissingArguments
	}

	if len(args) > 1 {
		return "", 0, errExtraArguments
	}

	if args[0] == "" {
		return "", 0, errInvalidProcessGuid
	}

	index := -1
	if indexFlag != "" {
		var err error
		index, err = strconv.Atoi(indexFlag)
		if err != nil || index < 0 {
			return "", 0, errInvalidIndex
		}
	}

	return args[0], index, nil
}

func ActualLRPsForGuid(stdout, stderr io.Writer, bbsClient bbs.Client, processGuid string, index int) error {
	logger := globalLogger.Session("actual-lrps-for-guid")

	encoder := json.NewEncoder(stdout)
	if index < 0 {
		actualLRPs, err := bbsClient.ActualLRPsByProcessGuid(logger, processGuid)
		if err != nil {
			return err
		}

		for _, lrp := range actualLRPs {
			err = encoder.Encode(lrp)
			if err != nil {
				logger.Error("failed-to-marshal", err)
			}
		}

		return nil
	} else {
		actualLRP, err := bbsClient.ActualLRPByProcessGuidAndIndex(logger, processGuid, index)
		if err != nil {
			return err
		}

		return encoder.Encode(actualLRP)
	}
}
