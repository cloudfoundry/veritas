package commands

import (
	"errors"
	"fmt"
	"os"

	"net/url"

	"code.cloudfoundry.org/bbs"
	"code.cloudfoundry.org/lager"
	"github.com/spf13/cobra"
)

var logger = lager.NewLogger("cfdot")

var RootCmd = &cobra.Command{
	Use:   "cfdot",
	Short: "Diego operator tooling",
	Long:  "A command-line tool to interact with a Cloud Foundry Diego deployment",
}

var (
	bbsURL            string
	bbsCACertFile     string
	bbsSkipCertVerify bool
	bbsCertFile       string
	bbsKeyFile        string
)

func addBBSFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&bbsURL, "bbsURL", "", "", "URL of BBS server to target, can also be specified with BBS_URL environment variable")
	cmd.Flags().StringVarP(&bbsCACertFile, "bbsCACertFile", "", "", "Path to CA file used to verify the BBS server")
	cmd.Flags().BoolVarP(&bbsSkipCertVerify, "bbsSkipCertVerify", "", false, "If set to true, do not verify the BBS server cert")
	cmd.Flags().StringVarP(&bbsCertFile, "bbsCertFile", "", "", "Path to cert file for the cfdot client to preset to the BBS for mutual TLS auth")
	cmd.Flags().StringVarP(&bbsKeyFile, "bbsKeyFile", "", "", "Path to the cfdot client key file used in mutual TLS auth")
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		if bbsURL == "" {
			bbsURL = os.Getenv("BBS_URL")
		}
		if bbsCACertFile == "" {
			bbsCACertFile = os.Getenv("BBS_CA_CERT_FILE")
		}
		if bbsSkipCertVerify == false {
			bbsCACertFile = os.Getenv("BBS_SKIP_CERT_VERIFY")
		}

		if bbsURL == "" {
			reportErr(cmd, errors.New(
				"BBS URL not set. Please specify one with the '--bbsURL' flag or the "+
					"'BBS_URL' environment variable.",
			), 3)
		} else if parsedURL, err := url.Parse(bbsURL); err != nil {
			reportErr(cmd, errors.New(fmt.Sprintf(
				"The value '%s' is not a valid BBS URL. Please specify one with the "+
					"'--bbsURL' flag or the 'BBS_URL' environment variable.",
				bbsURL,
			)), 3)
		} else if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
			reportErr(cmd, errors.New(fmt.Sprintf(
				"The URL '%s' does not have an 'http' or 'https' scheme. Please "+
					"specify one with the '--bbsURL' flag or the 'BBS_URL' environment "+
					"variable.",
				bbsURL,
			)), 3)
		}

	}
}

func newBBSClient(cmd *cobra.Command) bbs.Client {
	return bbs.NewClient(bbsURL)
}

func reportErr(cmd *cobra.Command, err error, exitCode int) {
	cmd.SetOutput(cmd.OutOrStderr())
	fmt.Fprintf(cmd.OutOrStderr(), "error: %s\n\n", err.Error())
	cmd.Help()
	os.Exit(exitCode)
}
