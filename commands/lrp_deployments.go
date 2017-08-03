package commands

import (
	"encoding/json"
	"io"

	"code.cloudfoundry.org/bbs"
	"github.com/spf13/cobra"
)

var lrpDeploymentsCmd = &cobra.Command{
	Use:   "lrp-deployments [[ids...]]",
	Short: "List all LRP Deployments, optionally filtering on specific deployment IDs",
	Long:  "Returns a list of all LRPDeployment objects. If a list of deployment IDs is passed, only deployments matching one of those IDs will be returned.",
	RunE:  lrpDeployments,
}

func init() {
	AddBBSFlags(lrpDeploymentsCmd)
	RootCmd.AddCommand(lrpDeploymentsCmd)
}

func lrpDeployments(cmd *cobra.Command, args []string) error {
	bbsClient, err := newBBSClient(cmd)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	err = LRPDeployments(cmd.OutOrStdout(), cmd.OutOrStderr(), bbsClient, args)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	return nil
}

func LRPDeployments(stdout, stderr io.Writer, bbsClient bbs.Client, args []string) error {
	logger := globalLogger.Session("lrp-deployments")

	lrpDeployments, err := bbsClient.LRPDeployments(logger, args)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(stdout)
	for _, lrp := range lrpDeployments {
		err = encoder.Encode(lrp)
		if err != nil {
			logger.Error("failed-to-marshal", err)
		}
	}

	return nil
}
