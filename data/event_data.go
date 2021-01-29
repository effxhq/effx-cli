package data

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	effx_api "github.com/effxhq/effx-api-v2/generated/go/client"
)

type EffxEvent struct {
	Payload *effx_api.CreateEventPayload
}

// Sync Updates the Effx API with yaml file contents
func (y EffxEvent) SendEvent(apiKey string) error {
	log.Printf("Sending event payload %+v\n", y.Payload)

	body, _ := json.Marshal(y.Payload)

	url := generateURL()
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
