package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"code.cloudfoundry.org/bbs"
	"code.cloudfoundry.org/bbs/models"
	"github.com/spf13/cobra"
)

var createLRPDeploymentCmd = &cobra.Command{
	Use:   "create-lrp-deployment (SPEC|@FILE)",
	Short: "Create a desired LRP Deployment",
	Long:  "Create a desired LRP Deployment from the given spec. Spec can either be json encoded desired-lrp, e.g. '{\"process_guid\":\"some-guid\"}' or a file containing json encoded lrp-deployment-definition, e.g. @/path/to/spec/file",
	RunE:  createLRPDeployment,
}

func init() {
	AddBBSFlags(createLRPDeploymentCmd)
	RootCmd.AddCommand(createLRPDeploymentCmd)
}

func createLRPDeployment(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return NewCFDotValidationError(cmd, fmt.Errorf("missing spec argument"))
	}

	spec, err := ValidateCreateLRPDeploymentArguments(args)
	if err != nil {
		return NewCFDotValidationError(cmd, err)
	}

	bbsClient, err := newBBSClient(cmd)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	err = CreateLRPDeployment(cmd.OutOrStdout(), cmd.OutOrStderr(), bbsClient, spec)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	return nil
}

func ValidateCreateLRPDeploymentArguments(args []string) ([]byte, error) {
	var lrpDeploymentDefinition *models.LRPDeploymentCreation
	var err error
	var spec []byte
	argValue := args[0]
	if strings.HasPrefix(argValue, "@") {
		_, err := os.Stat(argValue[1:])
		if err != nil {
			return nil, err
		}
		spec, err = ioutil.ReadFile(argValue[1:])
		if err != nil {
			return nil, err
		}

	} else {
		spec = []byte(argValue)
	}
	err = json.Unmarshal([]byte(spec), &lrpDeploymentDefinition)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Invalid JSON: %s", err.Error()))
	}
	return spec, nil
}

func CreateLRPDeployment(stdout, stderr io.Writer, bbsClient bbs.Client, spec []byte) error {
	logger := globalLogger.Session("create-desired-lrp")

	var lrpDeploymentDefinition *models.LRPDeploymentCreation
	err := json.Unmarshal(spec, &lrpDeploymentDefinition)
	if err != nil {
		return err
	}
	err = bbsClient.CreateLRPDeployment(logger, lrpDeploymentDefinition)
	if err != nil {
		return err
	}

	return nil
}
