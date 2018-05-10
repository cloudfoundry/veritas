package commands

import (
	"encoding/json"
	"io"

	"code.cloudfoundry.org/bbs"
	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/cfdot/commands/helpers"
	"github.com/spf13/cobra"
)

var (
	lrpEventsR1CellIdFlag string
)

var lrpEventsR1Cmd = &cobra.Command{
	Use:   "lrp-events-r1",
	Short: "Subscribe to BBS LRP events",
	Long:  "Subscribe to BBS LRP events",
	RunE:  lrpEventsR1,
}

// type LRPEvent struct {
// 	Type string      `json:"type"`
// 	Data interface{} `json:"data"`
// }
//
func init() {
	AddBBSFlags(lrpEventsR1Cmd)

	lrpEventsR1Cmd.Flags().StringVarP(&lrpEventsR1CellIdFlag, "cell-id", "c", "", "retrieve only events for the given cell id")

	RootCmd.AddCommand(lrpEventsR1Cmd)
}

func lrpEventsR1(cmd *cobra.Command, args []string) error {
	err := validateLRPEventsArguments(args)
	if err != nil {
		return NewCFDotValidationError(cmd, err)
	}

	bbsClient, err := helpers.NewBBSClient(cmd, Config)
	if err != nil {
		return NewCFDotError(cmd, err)
	}

	err = LRPEventsR1(cmd.OutOrStdout(), cmd.OutOrStderr(), bbsClient, lrpEventsCellIdFlag)
	if err != nil {
		return NewCFDotError(cmd, err)
	}
	return nil
}

func validateLRPEventsR1Arguments(args []string) error {
	if len(args) > 0 {
		return errExtraArguments
	}
	return nil
}

func LRPEventsR1(stdout, stderr io.Writer, bbsClient bbs.Client, cellID string) error {
	logger := globalLogger.Session("lrp-events")

	es, err := bbsClient.SubscribeToEventsR1ByCellID(logger, cellID)
	if err != nil {
		return models.ConvertError(err)
	}
	defer es.Close()
	encoder := json.NewEncoder(stdout)

	var lrpEvent LRPEvent
	for {
		event, err := es.Next()
		switch err {
		case nil:
			lrpEvent.Type = event.EventType()
			lrpEvent.Data = event
			err = encoder.Encode(lrpEvent)
			if err != nil {
				logger.Error("failed-to-marshal", err)
			}
		case io.EOF:
			return nil
		default:
			return err
		}
	}
}
