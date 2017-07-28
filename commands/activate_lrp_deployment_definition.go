package commands

import (
	"errors"
	"io"

	"code.cloudfoundry.org/bbs"
	"github.com/spf13/cobra"
)

var activateLRPDeploymentDefinitionCmd = &cobra.Command{
	Use:   "activate-lrp-definition PROCESS_GUID DEFINITION_ID",
	Short: "Activate the given definition id",
	Long:  "Activate the given definition id of the lrp deployment",
	RunE:  activateLRPDeploymentDefinition,
}

func init() {
	AddBBSFlags(activateLRPDeploymentDefinitionCmd)
	RootCmd.AddCommand(activateLRPDeploymentDefinitionCmd)
}

func activateLRPDeploymentDefinition(cmd *cobra.Command, args []string) error {
	processGuid, definitionID, err := ValidateActivateLRPDeploymentDefinitionArguments(args)
	if err != nil {
		return NewCFDotValidationError(cmd, err)
	}

	bbsClient, err := newBBSClient(cmd)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	err = ActivateLRPDeploymentDefinition(cmd.OutOrStdout(), cmd.OutOrStderr(), bbsClient, processGuid, definitionID)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	return nil
}

func ValidateActivateLRPDeploymentDefinitionArguments(args []string) (string, string, error) {
	if len(args) == 0 {
		return "", "", errMissingArguments
	}

	if len(args) > 2 {
		return "", "", errExtraArguments
	}

	if args[0] == "" {
		return "", "", errInvalidProcessGuid
	}

	if args[1] == "" {
		return "", "", errors.New("definition ID should be non empty string")

	}

	return args[0], args[1], nil
}

func ActivateLRPDeploymentDefinition(stdout, stderr io.Writer, bbsClient bbs.Client, processGuid, definitionID string) error {
	logger := globalLogger.Session("activate-lrp-definition")

	err := bbsClient.ActivateLRPDefinition(logger, processGuid, definitionID)
	if err != nil {
		return err
	}

	return nil
}
