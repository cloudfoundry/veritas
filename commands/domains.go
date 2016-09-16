package commands

import (
	"encoding/json"
	"io"

	"code.cloudfoundry.org/bbs"

	"github.com/spf13/cobra"
)

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "List domains",
	Long:  "List fresh domains from the BBS",
	RunE:  domains,
}

func init() {
	AddBBSFlags(domainsCmd)
	domainsCmd.PreRunE = BBSPrehook
	RootCmd.AddCommand(domainsCmd)
}

func domains(cmd *cobra.Command, args []string) error {
	var err error
	var bbsClient bbs.Client

	bbsClient, err = newBBSClient(cmd)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	err = Domains(cmd.OutOrStdout(), cmd.OutOrStderr(), bbsClient, args)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	return nil
}

func Domains(stdout, stderr io.Writer, bbsClient bbs.Client, args []string) error {
	logger := globalLogger.Session("domains")

	encoder := json.NewEncoder(stdout)
	domains, err := bbsClient.Domains(logger)
	if err != nil {
		return err
	}

	for _, domain := range domains {
		encoder.Encode(domain)
	}

	return nil
}