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

var updateLRPDeploymentCmd = &cobra.Command{
	Use:   "update-lrp-deployment (SPEC|@FILE)",
	Short: "update a desired LRP Deployment",
	Long:  "update a desired LRP Deployment from the given spec. Spec can either be json encoded desired-lrp, e.g. '{\"process_guid\":\"some-guid\"}' or a file containing json encoded desired-lrp, e.g. @/path/to/spec/file",
	RunE:  updateLRPDeployment,
}

func init() {
	AddBBSFlags(updateLRPDeploymentCmd)
	RootCmd.AddCommand(updateLRPDeploymentCmd)
}

func updateLRPDeployment(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return NewCFDotValidationError(cmd, fmt.Errorf("Missing arguments"))
	}

	processGuid, spec, err := ValidateUpdateLRPDeploymentArguments(args)
	if err != nil {
		return NewCFDotValidationError(cmd, err)
	}

	bbsClient, err := newBBSClient(cmd)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	err = UpdateLRPDeployment(cmd.OutOrStdout(), cmd.OutOrStderr(), bbsClient, processGuid, spec)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	return nil
}

func ValidateUpdateLRPDeploymentArguments(args []string) (string, []byte, error) {
	var lrpDeploymentUpdate *models.LRPDeploymentUpdate
	var err error
	var spec []byte
	processGuid := args[0]
	argValue := args[1]
	if strings.HasPrefix(argValue, "@") {
		_, err := os.Stat(argValue[1:])
		if err != nil {
			println(err.Error())
			return "", nil, err
		}
		spec, err = ioutil.ReadFile(argValue[1:])
		if err != nil {
			return "", nil, err
		}

	} else {
		spec = []byte(argValue)
	}
	err = json.Unmarshal([]byte(spec), &lrpDeploymentUpdate)
	if err != nil {
		return "", nil, errors.New(fmt.Sprintf("Invalid JSON: %s", err.Error()))
	}
	return processGuid, spec, nil
}

func UpdateLRPDeployment(stdout, stderr io.Writer, bbsClient bbs.Client, processGuid string, spec []byte) error {
	logger := globalLogger.Session("update-desired-lrp")

	var lrpDeploymentUpdate *models.LRPDeploymentUpdate
	err := json.Unmarshal(spec, &lrpDeploymentUpdate)
	if err != nil {
		return err
	}
	err = bbsClient.UpdateLRPDeployment(logger, processGuid, lrpDeploymentUpdate)
	if err != nil {
		return err
	}

	return nil
}
