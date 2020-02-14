package event

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/effxhq/effx-go/data"
	"github.com/effxhq/effx-go/internal/client/http"
	"github.com/effxhq/effx-go/internal/validator"
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

		tags := []*data.Tag{}
		hashtags := []string{}

		if tagsString != "" {
			tagsStringNoSpace := strings.Join(strings.Fields(tagsString), "")
			splitTagsString := strings.Split(tagsStringNoSpace, ",")

			for _, splitTag := range splitTagsString {
				splitTagString := strings.Split(splitTag, ":")

				if len(splitTagString) == 2 {
					tags = append(tags, &data.Tag{Key: splitTagString[0], Value: splitTagString[1]})
				} else {
					log.Fatalf("found invalid tag: %s", splitTag)
				}
			}
		}

		if hashtagsString != "" {
			hashtagsStringNoSpace := strings.Join(strings.Fields(hashtagsString), "")
			hashtags = strings.Split(hashtagsStringNoSpace, ",")
		}

		c := http.New(apiKeyString)

		object := &data.Data{
			Tap: &data.Tap{
				ProducedAtTimeMilliseconds: time.Now().UnixNano() / 1e6,
				Name:                       nameString,
				Description:                descriptionString,
				Tags:                       tags,
				Hashtags:                   hashtags,
				Integration: &data.Integration{
					Name:    integrationNameString,
					Version: integrationVersionString,
				},
				Service: &data.Service{
					Name: serviceNameString,
				},
			},
		}

		if userEmailString != "" {
			object.Tap.User = &data.User{
				Email: userEmailString,
			}
		}

		if err := validator.ValidateObject(object); err != nil {
			log.Fatal(err.Error())
		}

		if err := c.Synchronize(object); err != nil {
			log.Fatal(err.Error())
		}
	},
}
