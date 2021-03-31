package data

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	effx_api "github.com/effxhq/effx-api-v2-go/client"
)

type EffxEvent struct {
	Payload *effx_api.CreateEventPayload
}

// GetIntegrationName uses the environment to generate an effx integration name
func GetIntegrationName() string {
	integrations := map[string]string{
		"gitlab":         "GITLAB_CI",
		"github_actions": "GITHUB_ACTIONS",
		"circleci":       "CIRCLECI",
		"semaphore":      "SEMAPHORE",
	}

	for integration, envVariable := range integrations {
		if os.Getenv(envVariable) != "" {
			return integration
		}
	}

	return ""
}

// Sync Updates the Effx API with yaml file contents
func (y EffxEvent) SendEvent(apiKey string) error {
	body, _ := json.Marshal(y.Payload)

	log.Printf("Sending event payload %+v\n", string(body))

	url := GenerateUrl()
	url.Path = "v2/events"

	request, _ := http.NewRequest("POST", url.String(), bytes.NewReader(body))
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-effx-api-key", apiKey)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkForErrors(resp)
}
