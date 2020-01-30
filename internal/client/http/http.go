package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/effxhq/effx-go/data"
	"github.com/effxhq/effx-go/internal/client"
	"gopkg.in/go-playground/validator.v9"
)

const serviceIngestHost = "https://ingest.effx.io/v1/service"

type httpClient struct {
	apiKey   string
	validate *validator.Validate
}

func New(apiKey string) client.Client {
	return &httpClient{
		apiKey:   apiKey,
		validate: validator.New(),
	}
}

func (h *httpClient) Synchronize(object *data.Data) error {
	if err := h.validate.Struct(object); err != nil {
		return err
	}

	if object.Service != nil {
		h.synchronizeService(object)
	}

	return nil
}

func (h *httpClient) synchronizeService(object *data.Data) error {
	var (
		jsonPayload *bytes.Buffer
		err         error
		req         *http.Request
	)

	if object.Service == nil {
		return errors.New("service key is not found")
	}

	jsonPayload = new(bytes.Buffer)
	json.NewEncoder(jsonPayload).Encode(object.Service)

	req, err = http.NewRequest(serviceIngestHost, "application/json", jsonPayload)
	req.Header.Set("X-Effx-Api-Key", h.apiKey)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
