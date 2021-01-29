package event

import (
	"errors"
	"os"

	"github.com/effxhq/effx-cli/internal/parser"
	"github.com/spf13/cobra"
)

const effxApiKeyName = "EFFX_API_KEY"

var (
	apiKeyString string
	result       *parser.EventPayload = &parser.EventPayload{}
)

func init() {
	EventCmd.PersistentFlags().StringVarP(&apiKeyString, "key", "k", "", "your effx api key. alternatively, you can use env var EFFX_API_KEY")
	EventCmd.PersistentFlags().StringVarP(&result.Title, "title", "tt", "", "name of the event")
	EventCmd.PersistentFlags().StringVarP(&result.Message, "message", "msg", "", "message to describe the event")
	EventCmd.PersistentFlags().StringVarP(&result.ServiceName, "service", "sn", "", "service name the event is associated with")
	EventCmd.PersistentFlags().StringVarP(&result.Tags, "tags", "", "tg", "tags in the format of k:v,k1:v1 . use commas to separatewe tags")
	EventCmd.PersistentFlags().StringVarP(&result.Actions, "actions", "act", "", "actions in the format of <level>:<name>:<url>")
	EventCmd.PersistentFlags().IntVarP(&result.ProducedAtTimeMS, "produced_at_time", "timestamp", 0, "optional time the event was created at. format is epoch milliseconds. default is current time")
}

var EventCmd = &cobra.Command{
	Use:   "event",
	Short: "send events to the effx api",
	Long:  `send events to the effx api`,
	Args: func(cmd *cobra.Command, args []string) error {
		if apiKeyString == "" {
			if apiKeyString = os.Getenv(effxApiKeyName); apiKeyString == "" {
				return errors.New("api key is required")
			}
		}

		if result.Title == "" {
			return errors.New("--title <title> is required")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := parser.ProcessEvent(result)

		return payload.SendEvent(apiKeyString)
	},
}
