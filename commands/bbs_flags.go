package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"code.cloudfoundry.org/cfdot/commands/helpers"

	"net/url"

	"github.com/spf13/cobra"
)

var (
	BBSClientConfig helpers.BBSClientConfig
	bbsPreHooks     = []func(cmd *cobra.Command, args []string) error{}
)

// errors
var (
	errMissingBBSUrl = errors.New("BBS URL not set. Please specify one with the '--bbsURL' flag or the 'BBS_URL' environment variable.")
)

func AddBBSFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&BBSClientConfig.Url, "bbsURL", "", "URL of BBS server to target [environment variable equivalent: BBS_URL]")
	cmd.Flags().BoolVar(&BBSClientConfig.SkipCertVerify, "bbsSkipCertVerify", false, "when set to true, skips all SSL/TLS certificate verification [environment variable equivalent: BBS_SKIP_CERT_VERIFY]")
	cmd.Flags().StringVar(&BBSClientConfig.CertFile, "bbsCertFile", "", "path to the TLS client certificate to use during mutual-auth TLS [environment variable equivalent: BBS_CERT_FILE]")
	cmd.Flags().StringVar(&BBSClientConfig.KeyFile, "bbsKeyFile", "", "path to the TLS client private key file to use during mutual-auth TLS [environment variable equivalent: BBS_KEY_FILE]")
	cmd.Flags().StringVar(&BBSClientConfig.CACertFile, "bbsCACertFile", "", "path the Certificate Authority (CA) file to use when verifying TLS keypairs [environment variable equivalent: BBS_CA_CERT_FILE]")
	bbsPreHooks = append(bbsPreHooks, cmd.PreRunE)
	cmd.PreRunE = BBSPrehook
}

func BBSPrehook(cmd *cobra.Command, args []string) error {
	var err error
	for _, f := range bbsPreHooks {
		if f == nil {
			continue
		}
		err = f(cmd, args)
		if err != nil {
			return err
		}
	}
	BBSClientConfig.TLSConfig.Merge(Config)
	return setBBSFlags(cmd, args)
}

func setBBSFlags(cmd *cobra.Command, args []string) error {
	var err, returnErr error

	if BBSClientConfig.Url == "" {
		BBSClientConfig.Url = os.Getenv("BBS_URL")
	}

	// Only look at the environment variable if the flag has not been set.
	if !cmd.Flags().Lookup("bbsSkipCertVerify").Changed && os.Getenv("BBS_SKIP_CERT_VERIFY") != "" {
		BBSClientConfig.SkipCertVerify, err = strconv.ParseBool(os.Getenv("BBS_SKIP_CERT_VERIFY"))
		if err != nil {
			returnErr = NewCFDotValidationError(
				cmd,
				fmt.Errorf(
					"The value '%s' is not a valid value for BBS_SKIP_CERT_VERIFY. Please specify one of the following valid boolean values: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False",
					os.Getenv("BBS_SKIP_CERT_VERIFY")),
			)
			return returnErr
		}
	}

	if BBSClientConfig.CertFile == "" {
		BBSClientConfig.CertFile = os.Getenv("BBS_CERT_FILE")
	}

	if BBSClientConfig.KeyFile == "" {
		BBSClientConfig.KeyFile = os.Getenv("BBS_KEY_FILE")
	}

	if BBSClientConfig.CACertFile == "" {
		BBSClientConfig.CACertFile = os.Getenv("BBS_CA_CERT_FILE")
	}

	if BBSClientConfig.Url == "" {
		returnErr = NewCFDotValidationError(cmd, errMissingBBSUrl)
		return returnErr
	}

	var parsedURL *url.URL
	if parsedURL, err = url.Parse(BBSClientConfig.Url); err != nil {
		returnErr = NewCFDotValidationError(
			cmd,
			fmt.Errorf(
				"The value '%s' is not a valid BBS URL. Please specify one with the '--bbsURL' flag or the 'BBS_URL' environment variable.",
				BBSClientConfig.Url),
		)
		return returnErr
	}

	if parsedURL.Scheme == "https" {
		if !BBSClientConfig.SkipCertVerify {
			if BBSClientConfig.CACertFile == "" {
				returnErr = NewCFDotValidationError(cmd, errMissingCACertFile)
				return returnErr
			}

			err := validateReadableFile(cmd, BBSClientConfig.CACertFile, "CA cert")
			if err != nil {
				return err
			}
		}

		if (BBSClientConfig.KeyFile == "") || (BBSClientConfig.CertFile == "") {
			returnErr = NewCFDotValidationError(cmd, errMissingClientCertAndKeyFiles)
			return returnErr
		}

		if BBSClientConfig.KeyFile != "" {
			err := validateReadableFile(cmd, BBSClientConfig.KeyFile, "key")

			if err != nil {
				return err
			}
		}

		if BBSClientConfig.CertFile != "" {
			err := validateReadableFile(cmd, BBSClientConfig.CertFile, "cert")

			if err != nil {
				return err
			}
		}

		return nil
	}

	if parsedURL.Scheme != "http" {
		returnErr = NewCFDotValidationError(
			cmd,
			fmt.Errorf(
				"The URL '%s' does not have an 'http' or 'https' scheme. Please "+
					"specify one with the '--bbsURL' flag or the 'BBS_URL' environment "+
					"variable.", BBSClientConfig.Url),
		)
		return returnErr
	}

	return nil
}

func validateReadableFile(cmd *cobra.Command, filename, filetype string) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {

		return NewCFDotValidationError(
			cmd,
			fmt.Errorf("%s file '%s' doesn't exist or is not readable: %s", filetype, filename, err.Error()),
		)
	}

	return nil
}
