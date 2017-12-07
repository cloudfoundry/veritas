package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"code.cloudfoundry.org/cfdot/commands/helpers"
	"github.com/spf13/cobra"
)

var (
	LocketClientConfig helpers.LocketClientConfig
	locketPreHooks     = []func(cmd *cobra.Command, args []string) error{}
)

// errors
var (
	errMissingLocketUrl             = errors.New("Locket API Location not set. Please specify one with the '--locketAPILocation' flag or the 'LOCKET_API_LOCATION' environment variable.")
	errMissingLocketCACertFile      = errors.New("--locketCACertFile must be specified if --locketSkipCertVerify is not set")
	errMissingLocketCertAndKeyFiles = errors.New("--locketCertFile and --locketKeyFile must both be specified for TLS connections.")
)

func AddLocketFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&LocketClientConfig.ApiLocation, "locketAPILocation", "", "Hostname:Port of Locket server to target [environment variable equivalent: LOCKET_API_LOCATION]")
	cmd.Flags().BoolVar(&LocketClientConfig.SkipCertVerify, "locketSkipCertVerify", false, "when set to true, skips all SSL/TLS certificate verification [environment variable equivalent: LOCKET_SKIP_CERT_VERIFY]")
	cmd.Flags().StringVar(&LocketClientConfig.CertFile, "locketCertFile", "", "path to the TLS client certificate to use during mutual-auth TLS [environment variable equivalent: LOCKET_CERT_FILE]")
	cmd.Flags().StringVar(&LocketClientConfig.KeyFile, "locketKeyFile", "", "path to the TLS client private key file to use during mutual-auth TLS [environment variable equivalent: LOCKET_KEY_FILE]")
	cmd.Flags().StringVar(&LocketClientConfig.CACertFile, "locketCACertFile", "", "path the Certificate Authority (CA) file to use when verifying TLS keypairs [environment variable equivalent: LOCKET_CA_CERT_FILE]")
	bbsPreHooks = append(bbsPreHooks, cmd.PreRunE)
	cmd.PreRunE = LocketPrehook
}

func LocketPrehook(cmd *cobra.Command, args []string) error {
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
	LocketClientConfig.TLSConfig.Merge(Config)
	return setLocketFlags(cmd, args)
}

func setLocketFlags(cmd *cobra.Command, args []string) error {
	var err, returnErr error
	if LocketClientConfig.ApiLocation == "" {
		LocketClientConfig.ApiLocation = os.Getenv("LOCKET_API_LOCATION")
	}

	// Only look at the environment variable if the flag has not been set.
	if !cmd.Flags().Lookup("locketSkipCertVerify").Changed && os.Getenv("LOCKET_SKIP_CERT_VERIFY") != "" {
		LocketClientConfig.SkipCertVerify, err = strconv.ParseBool(os.Getenv("LOCKET_SKIP_CERT_VERIFY"))
		if err != nil {
			returnErr = NewCFDotValidationError(
				cmd,
				fmt.Errorf(
					"The value '%s' is not a valid value for LOCKET_SKIP_CERT_VERIFY. Please specify one of the following valid boolean values: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False",
					os.Getenv("LOCKET_SKIP_CERT_VERIFY")),
			)
			return returnErr
		}
	}

	if LocketClientConfig.CertFile == "" {
		LocketClientConfig.CertFile = os.Getenv("LOCKET_CERT_FILE")
	}

	if LocketClientConfig.KeyFile == "" {
		LocketClientConfig.KeyFile = os.Getenv("LOCKET_KEY_FILE")
	}

	if LocketClientConfig.CACertFile == "" {
		LocketClientConfig.CACertFile = os.Getenv("LOCKET_CA_CERT_FILE")
	}

	if LocketClientConfig.ApiLocation == "" {
		returnErr = NewCFDotValidationError(cmd, errMissingLocketUrl)
		return returnErr
	}

	if !LocketClientConfig.SkipCertVerify {
		if LocketClientConfig.CACertFile == "" {
			returnErr = NewCFDotValidationError(cmd, errMissingCACertFile)
			return returnErr
		}

		err := validateReadableFile(cmd, LocketClientConfig.CACertFile, "CA cert")

		if err != nil {
			return err
		}
	}

	if (LocketClientConfig.KeyFile == "") || (LocketClientConfig.CertFile == "") {
		returnErr = NewCFDotValidationError(cmd, errMissingClientCertAndKeyFiles)
		return returnErr
	}

	if LocketClientConfig.KeyFile != "" {
		err := validateReadableFile(cmd, LocketClientConfig.KeyFile, "key")

		if err != nil {
			return err
		}
	}

	if LocketClientConfig.CertFile != "" {
		err := validateReadableFile(cmd, LocketClientConfig.CertFile, "cert")

		if err != nil {
			return err
		}
	}

	return nil
}
