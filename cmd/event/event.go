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
	result       *parser.EventPayload
)

func init() {
	eventCmd.PersistentFlags().StringVarP(&apiKeyString, "key", "k", "", "your effx api key. alternatively, you can use env var EFFX_API_KEY")
	eventCmd.PersistentFlags().StringVarP(&result.Name, "title", "", "", "name of the event")
	eventCmd.PersistentFlags().StringVarP(&result.Title, "message", "", "", "message to describe the event")
	eventCmd.PersistentFlags().StringVarP(&result.ServiceName, "service", "", "", "service name the event is associated with")
	eventCmd.PersistentFlags().StringVarP(&result.Tags, "tags", "", "", "tags in the format of k:v,k1:v1 . use commas to separatewe tags")
	eventCmd.PersistentFlags().StringVarP(&result.Actions, "actions", "", "", "actions in the format of <level>:<name>:<url>")
	eventCmd.PersistentFlags().IntVarP(&result.ProducedAtTimeMS, "produced_at_time", "", 0, "optional time the event was created at. format is epoch milliseconds. default is current time")
}

var eventCmd = &cobra.Command{
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
		payload := parser.ProcessEvent(&parser.EventPayload{
			Name:             result.Name,
			Message:          result.Message,
			ServiceName:      result.ServiceName,
			Tags:             result.Tags,
			Actions:          result.Actions,
			ProducedAtTimeMS: result.ProducedAtTimeMS,
		})

		return payload.SendEvent(apiKeyString)
	},
}
