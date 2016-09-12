package desired_lrp_scheduling_infos_command

import (
	"encoding/json"
	"io"

	. "code.cloudfoundry.org/cfdot/commands"

	"code.cloudfoundry.org/bbs"
	"code.cloudfoundry.org/bbs/models"

	"github.com/spf13/cobra"
)

var (
	domain string
)

var desiredLRPSchedulingInfosCmd = &cobra.Command{
	Use:   "desired-lrp-scheduling-infos",
	Short: "List desired LRP scheduling infos",
	Long:  "List desired LRP scheduling infos from the BBS",
	RunE:  desiredLRPSchedulingInfos,
}

func init() {
	AddBBSFlags(desiredLRPSchedulingInfosCmd)
	desiredLRPSchedulingInfosCmd.PreRunE = BBSPrehook
	desiredLRPSchedulingInfosCmd.Flags().StringVarP(&domain, "domain", "d", "", "retrieve only scheduling infos for the given domain")
	RootCmd.AddCommand(desiredLRPSchedulingInfosCmd)
}

func desiredLRPSchedulingInfos(cmd *cobra.Command, args []string) error {
	var err error
	var bbsClient bbs.Client

	bbsClient, err = NewBBSClient(cmd)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	err = DesiredLRPSchedulingInfos(cmd.OutOrStdout(), cmd.OutOrStderr(), bbsClient, args)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	return nil
}

func DesiredLRPSchedulingInfos(stdout, stderr io.Writer, bbsClient bbs.Client, args []string) error {
	logger := GlobalLogger.Session("desiredLRPSchedulingInfos")

	encoder := json.NewEncoder(stdout)
	desiredLRPFilter := models.DesiredLRPFilter{}

	if domain != "" {
		desiredLRPFilter.Domain = domain
	}

	desiredLRPSchedulingInfos, err := bbsClient.DesiredLRPSchedulingInfos(logger, desiredLRPFilter)
	if err != nil {
		return err
	}

	for _, desiredLRPSchedulingInfo := range desiredLRPSchedulingInfos {
		encoder.Encode(desiredLRPSchedulingInfo)
	}

	return nil
}
