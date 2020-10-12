package data

import (
	"context"
	"encoding/json"
	"log"

	effx_api "github.com/effxhq/effx-api-v2/generated/go/client"
	validate "github.com/effxhq/openapi3-validate"
)

// Store as JSON bytes
type ApiResourceContent struct {
	Content []byte
}

type ApiResourceMeta struct {
	Kind       string
	ApiVersion string
	Metadata   struct {
		Name string
	}
}

type ApiResource interface {
	GetMeta() ApiResourceMeta
	Lint() error
	Sync(apiKey string, isPost bool) error
}

func (c ApiResourceContent) GetMeta() ApiResourceMeta {
	var meta ApiResourceMeta
	_ = json.Unmarshal(c.Content, &meta)
	return meta
}

func (c ApiResourceContent) Lint() error {
	log.Printf("Linting %+v\n", c.GetMeta())
	return validate.ValidateComponent(c.Content)
}

func (c ApiResourceContent) Sync(apiKey string, isPost bool) error {
	log.Printf("Syncing %+v\n", c.GetMeta())
	cfg := effx_api.NewConfiguration()
	client := effx_api.NewAPIClient(cfg)

	serverIndex := 0
	if isPost {
		serverIndex = 1
	}
	ctx := context.WithValue(context.Background(), effx_api.ContextServerIndex, serverIndex)

	meta := c.GetMeta()
	switch meta.Kind {
	case "Service":
		var service effx_api.ServiceConfiguration
		if err := json.Unmarshal(c.Content, &service); err != nil {
			return err
		}
		request := client.ServicesApi.ServicesPut(ctx)
		request = request.ServiceConfiguration(service)
		request = request.XEffxApiKey(apiKey)
		_, err := request.Execute()
		if err != nil {
			return err
		}
	case "Team":
		var team effx_api.TeamConfiguration
		if err := json.Unmarshal(c.Content, &team); err != nil {
			return err
		}
		request := client.TeamsApi.TeamsPut(ctx)
		request = request.TeamConfiguration(team)
		request = request.XEffxApiKey(apiKey)
		_, err := request.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}
