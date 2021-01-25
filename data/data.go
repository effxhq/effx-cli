package data

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	effx_api "github.com/effxhq/effx-api-v2/generated/go/client"
	effx_v1_api "github.com/effxhq/effx-api/generated/go"
)

// EffxAPIHost Is the environment variable to override the API host
const EffxAPIHost = "EFFX_API_HOST"

// EffxYamlPattern is the regex pattern for yaml files
const EffxYamlPattern = "(.+\\.)?effx\\.ya?ml$"

var effxYamlRegex = regexp.MustCompile(EffxYamlPattern)

// EffxYaml provides a data structure and methods for interacting with effx yamls
type EffxYaml struct {
	FilePath string
}

func (y EffxYaml) isEffxYaml() bool {
	matched := effxYamlRegex.MatchString(y.FilePath)
	return matched
}

func (y EffxYaml) newConfig() (*effx_api.ConfigurationFile, error) {
	config := &effx_api.ConfigurationFile{}
	yamlFile, err := ioutil.ReadFile(y.FilePath)
	if err != nil {
		return nil, err
	}
	config.FileContents = string(yamlFile)
	config.SetAnnotations(map[string]string{
		"effx.io/source":    "effx-cli",
		"effx.io/file-path": y.FilePath,
	})

	return config, nil
}

// Lint Checks for syntax errors in the yaml file
func (y EffxYaml) Lint() error {
	log.Printf("Linting %+v\n", y.FilePath)

	ok := y.isEffxYaml()
	if !ok {
		errString := fmt.Sprintf("Not an Effx Yaml. %s must match pattern: %s", y.FilePath, EffxYamlPattern)
		return errors.New(errString)
	}

	config, err := y.newConfig()
	if err != nil {
		return nil
	}
	body, _ := json.Marshal(config)

	url := generateURL()
	url.Path = "v2/config/lint"

	resp, err := http.Post(url.String(), "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkForErrors(resp)
}

// Sync Updates the Effx API with yaml file contents
func (y EffxYaml) Sync(apiKey string) error {
	log.Printf("Syncing %+v\n", y.FilePath)

	config, err := y.newConfig()
	if err != nil {
		return nil
	}
	body, _ := json.Marshal(config)

	url := generateURL()
	url.Path = "v2/config"

	request, _ := http.NewRequest("PUT", url.String(), bytes.NewReader(body))
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-effx-api-key", apiKey)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkForErrors(resp)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func generateURL() *url.URL {
	u := url.URL{
		Scheme: "https",
		Host:   getEnv(EffxAPIHost, "api.effx.io"),
	}

	if strings.HasPrefix(u.Host, "localhost") {
		u.Scheme = "http"
	}

	return &u
}

func checkForErrors(response *http.Response) error {
	if response.StatusCode != 204 {
		var result map[string]interface{}
		_ = json.NewDecoder(response.Body).Decode(&result)
		return fmt.Errorf("%d: %s", response.StatusCode, result)
	}
	return nil
}

// Data for interacting with the v1 api
type V1Data struct {
	Service *effx_v1_api.ServicePayload
	Event   *effx_v1_api.EventPayload
	Team    *effx_v1_api.TeamPayload
	User    *effx_v1_api.UserPayload
}

func (d *V1Data) SendEvent(apiKey string) error {
	basePath := "https://api.effx.io/v1"

	if os.Getenv("EFFX_API_HOST") != "" {
		log.Printf("switching to use basePath %s", fmt.Sprintf("%s/v1", os.Getenv("EFFX_API_HOST")))
		basePath = fmt.Sprintf("%s/v1", os.Getenv("EFFX_API_HOST"))
	}

	client := effx_v1_api.NewAPIClient(&effx_v1_api.Configuration{
		BasePath:      basePath,
		DefaultHeader: make(map[string]string),
		UserAgent:     "Swagger-Codegen/1.0.0/go",
	})

	resp, err := client.EventsApi.EventsPut(
		context.Background(),
		apiKey,
		*d.Event,
		nil,
	)

	if err != nil {
		log.Printf("error %s", string(err.(effx_v1_api.GenericSwaggerError).Body()))
		log.Fatal(err.Error())
		return err
	}

	return checkForErrors(resp)
}
