package event

import (
	"errors"
	"os"

	// "github.com/effxhq/effx-cli/data"
	"github.com/effxhq/effx-cli/internal/parser"
	"github.com/spf13/cobra"
)

const effxApiKeyName = "EFFX_API_KEY"

var (
	apiKeyString string
	name         string
	title        string
	serviceName  string
	tags         string
	actions      string
)

func init() {
	eventCmd.PersistentFlags().StringVarP(&apiKeyString, "key", "k", "", "your effx api key. alternatively, you can use env var EFFX_API_KEY")
	eventCmd.PersistentFlags().StringVarP(&name, "title", "", "", "name of the event")
	eventCmd.PersistentFlags().StringVarP(&title, "message", "", "", "message to describe the event")
	eventCmd.PersistentFlags().StringVarP(&serviceName, "service", "", "", "service name the event is associated with")
	eventCmd.PersistentFlags().StringVarP(&tags, "tags", "", "", "tags in the format of k:v . use commas to separate tags")
	eventCmd.PersistentFlags().StringVarP(&actions, "actions", "", "", "actions in the format of <level>:<name>:<url>")
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

		if title == "" {
			return errors.New("--title <title> is required")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := parser.ProcessEvent(&parser.EventPayload{
			Name:        name,
			Message:     title,
			ServiceName: serviceName,
			Tags:        tags,
			Actions:     actions,
		})

		return payload.SendEvent(apiKeyString)
	},
}
