package commands

import (
	"io"

	"code.cloudfoundry.org/bbs"
	"github.com/spf13/cobra"
)

var deleteLRPDeploymentCmd = &cobra.Command{
	Use:   "delete-lrp-deployment PROCESS_GUID",
	Short: "Delete a desired LRP",
	Long:  "Delete a desired LRP with the given process guid.",
	RunE:  deleteLRPDeployment,
}

func init() {
	AddBBSFlags(deleteLRPDeploymentCmd)
	RootCmd.AddCommand(deleteLRPDeploymentCmd)
}

func deleteLRPDeployment(cmd *cobra.Command, args []string) error {
	processGuid, err := ValidateDeleteLRPDeploymentArguments(args)
	if err != nil {
		return NewCFDotValidationError(cmd, err)
	}

	bbsClient, err := newBBSClient(cmd)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	err = DeleteLRPDeployment(cmd.OutOrStdout(), cmd.OutOrStderr(), bbsClient, processGuid)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	return nil
}

func ValidateDeleteLRPDeploymentArguments(args []string) (string, error) {
	if len(args) == 0 {
		return "", errMissingArguments
	}

	if len(args) > 1 {
		return "", errExtraArguments
	}

	if args[0] == "" {
		return "", errInvalidProcessGuid
	}

	return args[0], nil
}

func DeleteLRPDeployment(stdout, stderr io.Writer, bbsClient bbs.Client, processGuid string) error {
	logger := globalLogger.Session("delete-lrp-deployment")

	err := bbsClient.DeleteLRPDeployment(logger, processGuid)
	if err != nil {
		return err
	}

	return nil
}
