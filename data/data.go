package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	effx_api "github.com/effxhq/effx-api-v2/generated/go/client"
	"github.com/effxhq/effx-cli/metadata"
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

func setMetadata(config *effx_api.ConfigurationFile, m *metadata.Result) *effx_api.ConfigurationFile {
	if m != nil {
		(*config.Annotations)["effx.io/inferred-tags"] = fmt.Sprintf("language,%s", m.Language)
		(*config.Tags)["language"] = m.Language
		(*config.Tags)[m.Language] = m.Version
	}

	return config
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

	result, err := metadata.InferMetadata(filepath.Dir(y.FilePath))
	if err != nil {
		log.Printf("Could not predict version %+v\n", err)
	}

	config = setMetadata(config, result)

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
