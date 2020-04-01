package event

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/antihax/optional"
	effx_api "github.com/effxhq/effx-api/generated/go"
	"github.com/effxhq/effx-cli/data"
	"github.com/spf13/cobra"
)

const effxApiKeyName = "EFFX_API_KEY"

var (
	apiKeyString             string
	nameString               string
	descriptionString        string
	serviceNameString        string
	integrationNameString    string
	integrationVersionString string
	isDryRun                 bool
	userEmailString          string
	tagsString               string
	hashtagsString           string
)

func Initialize() {
	EventCreateCmd.PersistentFlags().StringVarP(&apiKeyString, "key", "k", "", "your effx api key. alternatively, you can use env var EFFX_API_KEY")
	EventCreateCmd.PersistentFlags().StringVarP(&nameString, "name", "", "", "name of your event")
	EventCreateCmd.PersistentFlags().StringVarP(&descriptionString, "desc", "", "", "name of your event")
	EventCreateCmd.PersistentFlags().StringVarP(&serviceNameString, "service", "", "", "name of service")
	EventCreateCmd.PersistentFlags().StringVarP(&userEmailString, "user", "", "", "email for current user")
	EventCreateCmd.PersistentFlags().StringVarP(&integrationNameString, "integration_name", "", "", "name of integration")
	EventCreateCmd.PersistentFlags().StringVarP(&integrationVersionString, "integration_version", "", "", "version of integration")
	EventCreateCmd.PersistentFlags().StringVarP(&tagsString, "tags", "t", "", "tags in the format of k:v . use commas to separate tags")
	EventCreateCmd.PersistentFlags().StringVarP(&hashtagsString, "hashtags", "", "", "hashtags. use commas to separate hashtags")
	EventCreateCmd.PersistentFlags().BoolVarP(&isDryRun, "dry-run", "", false, "validate file(s)")
}

var EventCreateCmd = &cobra.Command{
	Use:   "event create",
	Short: "create an event via the effx api",
	Long:  `create an event via the effx api`,
	Run: func(cmd *cobra.Command, args []string) {
		if apiKeyString == "" {
			if apiKeyString = os.Getenv(effxApiKeyName); apiKeyString == "" {
				log.Fatal("api key is required")
			}
		}

		tags := []effx_api.TagPayload{}
		hashtags := []string{}

		if tagsString != "" {
			tagsStringNoSpace := strings.Join(strings.Fields(tagsString), "")
			splitTagsString := strings.Split(tagsStringNoSpace, ",")

			for _, splitTag := range splitTagsString {
				splitTagString := strings.Split(splitTag, ":")

				if len(splitTagString) == 2 {
					tags = append(tags, effx_api.TagPayload{Key: splitTagString[0], Value: splitTagString[1]})
				} else {
					log.Fatalf("found invalid tag: %s", splitTag)
				}
			}
		}

		if hashtagsString != "" {
			hashtagsStringNoSpace := strings.Join(strings.Fields(hashtagsString), "")
			hashtags = strings.Split(hashtagsStringNoSpace, ",")
		}

		object := &data.Data{
			Event: &effx_api.EventPayload{
				ProducedAtTimeMilliseconds: time.Now().UnixNano() / 1e6,
				Name:                       nameString,
				Description:                descriptionString,
				Tags:                       tags,
				Hashtags:                   hashtags,
				Integration: &effx_api.IntegrationPayload{
					Name:    integrationNameString,
					Version: integrationVersionString,
				},
				Service: &effx_api.EventServicePayload{
					Name: serviceNameString,
				},
			},
		}

		if userEmailString != "" {
			object.Event.User = &effx_api.EventUserPayload{
				Email: userEmailString,
			}
		}
		client := effx_api.NewAPIClient(effx_api.NewConfiguration())

		_, err := client.EventsApi.EventsPut(
			context.Background(),
			apiKeyString,
			*object.Event,
			&effx_api.EventsPutOpts{
				XEffxValidateOnly: optional.NewString(fmt.Sprintf("%v", isDryRun)),
			},
		)

		if err != nil {
			log.Printf("error %s", string(err.(effx_api.GenericSwaggerError).Body()))
			log.Fatal(err.Error())
		}
	},
}
