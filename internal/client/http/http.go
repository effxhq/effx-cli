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
		if err := h.synchronizeService(object); err != nil {
			return err
		}
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

	if err := json.NewEncoder(jsonPayload).Encode(object.Service); err != nil {
		return err
	}

	req, err = http.NewRequest("POST", serviceIngestHost, jsonPayload)

	if err != nil {
		return err
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("X-Effx-Api-Key", h.apiKey)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
