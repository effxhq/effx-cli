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
	"regexp"

	effx_api "github.com/effxhq/effx-api-v2/generated/go/client"
)

const EffxApiHost = "EFFX_API_HOST"
const EffxYamlPattern = "(.+\\.)?effx\\.ya?ml$"

var EffxYamlRegex = regexp.MustCompile(EffxYamlPattern)

type EffxYaml struct {
	FilePath string
}

func (y EffxYaml) isEffxYaml() bool {
	matched := EffxYamlRegex.MatchString(y.FilePath)
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

	url := generateUrl()
	url.Path = "v2/config/lint"

	resp, err := http.Post(url.String(), "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	logErrorMessages(resp)

	return nil
}

func (y EffxYaml) Sync(apiKey string) error {
	log.Printf("Syncing %+v\n", y.FilePath)

	config, err := y.newConfig()
	if err != nil {
		return nil
	}
	body, _ := json.Marshal(config)

	url := generateUrl()
	url.Path = "v2/config"

	request, _ := http.NewRequest("PUT", url.String(), bytes.NewReader(body))
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-effx-api-key", apiKey)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	logErrorMessages(resp)

	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func generateUrl() *url.URL {
	url := url.URL{
		Scheme: "https",
		Host:   getEnv(EffxApiHost, "api.effx.io"),
	}
	return &url
}

func logErrorMessages(response *http.Response) {
	if response.StatusCode != 204 {
		var result map[string]interface{}
		_ = json.NewDecoder(response.Body).Decode(&result)
		log.Println(result["message"])
	}
}
