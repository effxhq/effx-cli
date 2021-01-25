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
	apiKeyString    string
	name            string
	description     string
	serviceName     string
	integrationName string
	imageUrl        string
	email           string
	tags            string
	hashtags        string
)

func init() {
	eventCmd.PersistentFlags().StringVarP(&apiKeyString, "key", "k", "", "your effx api key. alternatively, you can use env var EFFX_API_KEY")
	eventCmd.PersistentFlags().StringVarP(&name, "name", "", "", "name of the event")
	eventCmd.PersistentFlags().StringVarP(&description, "desc", "", "", "text to describe the event")
	eventCmd.PersistentFlags().StringVarP(&serviceName, "service", "", "", "service name the event is associated with")
	eventCmd.PersistentFlags().StringVarP(&integrationName, "integration_name", "", "", "name of integration")
	eventCmd.PersistentFlags().StringVarP(&imageUrl, "image", "img", "", "image url for the event")
	eventCmd.PersistentFlags().StringVarP(&email, "email", "", "", "email for current user")
	eventCmd.PersistentFlags().StringVarP(&tags, "tags", "", "", "tags in the format of k:v . use commas to separate tags")
	eventCmd.PersistentFlags().StringVarP(&hashtags, "hashtags", "", "", "hashtags. use commas to separate hashtags")
}

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "event effx.yaml file(s) to the effx api",
	Long:  `event effx.yaml file(s) to the effx api`,
	Args: func(cmd *cobra.Command, args []string) error {
		if apiKeyString == "" {
			if apiKeyString = os.Getenv(effxApiKeyName); apiKeyString == "" {
				return errors.New("api key is required")
			}
		}

		if description == "" {
			return errors.New("--description <description> is required")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := parser.ProcessEvent(&parser.EventPayload{
			Name:            name,
			Description:     description,
			ServiceName:     serviceName,
			IntegrationName: integrationName,
			ImageUrl:        imageUrl,
			Email:           email,
			Tags:            tags,
			Hashtags:        hashtags,
		})

		return payload.SendEvent(apiKeyString)
	},
}
