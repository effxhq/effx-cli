package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/effxhq/effx-go/data"
	"github.com/effxhq/effx-go/internal/client"
	"gopkg.in/go-playground/validator.v9"
)

const (
	serviceIngestHost = "https://ingest.effx.io/v1/service"
	tapIngestHost     = "https://ingest.effx.io/v1/events"
)

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

	if object.Tap != nil {
		if err := h.synchronizeTap(object); err != nil {
			return err
		}
	}

	if object.Service != nil {
		if err := h.synchronizeService(object); err != nil {
			return err
		}
	}

	return nil
}

func (h *httpClient) synchronizeTap(object *data.Data) error {
	var (
		jsonPayload *bytes.Buffer
	)

	if object.Tap == nil {
		return errors.New("tap key is not found")
	}

	jsonPayload = new(bytes.Buffer)

	if err := json.NewEncoder(jsonPayload).Encode(object.Tap); err != nil {
		return err
	}

	return h.makeEffxRequest(tapIngestHost, jsonPayload)
}

func (h *httpClient) synchronizeService(object *data.Data) error {
	var (
		jsonPayload *bytes.Buffer
	)

	if object.Service == nil {
		return errors.New("service key is not found")
	}

	jsonPayload = new(bytes.Buffer)

	if err := json.NewEncoder(jsonPayload).Encode(object.Service); err != nil {
		return err
	}

	return h.makeEffxRequest(serviceIngestHost, jsonPayload)
}

func (h *httpClient) makeEffxRequest(url string, jsonPayload io.Reader) error {
	req, err := http.NewRequest("POST", url, jsonPayload)

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
