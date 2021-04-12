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

	effx_api "github.com/effxhq/effx-api-v2-go/client"
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

// given a absolute path (example: /runner/folder/root/data/effx.yaml)
// and given a workingDir (example: root)
// it will return data/effx.yaml
func parseRelativePathFromAbsolutePath(absoluteDir, workingDir string) string {
	// if working directory does not end with a slash, add it
	if strings.LastIndex(workingDir, "/") != len(workingDir)-1 {
		workingDir += "/"
	}

	res := strings.Split(absoluteDir, workingDir)
	if len(res) > 1 {
		return res[1]
	}
	return absoluteDir
}

// given an absolute path, this function will find the root directory
// (directory with a .git file)
func findGitRootDirectory(absoluteDir string) (string, error) {
	pathDir := filepath.Dir(absoluteDir)

	for pathDir != "" {
		pathDir = filepath.Clean(pathDir)
		files, err := ioutil.ReadDir(pathDir)

		if err != nil {
			return "", err
		}

		for _, file := range files {
			if !file.IsDir() {
				if strings.Contains(file.Name(), ".git") {
					return filepath.Base(pathDir), nil
				}
			}
		}

		pathDir = filepath.Join(pathDir, "..")
	}

	return absoluteDir, nil
}

// converts /runner/sandbox/root/effx.yaml
// to root/effx.yaml for a more helpful file source.
// it will try to get working directory from circleci,
// if cannot find it, it will attempt to find the root directory
// where git files are present.
func convertToRelativePath(absoluteDir string) (string, error) {
	var (
		workDir = ""
		err     error
	)

	if workDir == "" {
		absoluteDir, err = filepath.Abs(absoluteDir)
		if err != nil {
			return absoluteDir, err
		}
		workDir, err = findGitRootDirectory(absoluteDir)
		if err != nil {
			return absoluteDir, err
		}
	}

	return parseRelativePathFromAbsolutePath(absoluteDir, workDir), nil
}

func setMetadata(config *effx_api.ConfigurationFile, m *metadata.Result) *effx_api.ConfigurationFile {
	var (
		inferredTags = []string{}
	)

	if m == nil {
		return config
	}

	tags, ok := config.GetTagsOk()
	if !ok {
		tags = &map[string]string{}
	}

	annotations, ok := config.GetAnnotationsOk()
	if !ok {
		annotations = &map[string]string{}
	}

	if m.Language != "" {
		languageTag := "language"
		(*tags)[languageTag] = m.Language
		inferredTags = append(inferredTags, languageTag)
	}

	if m.Language != "" && m.Version != "" {
		langVersionTag := strings.ToLower(m.Language)
		(*tags)[langVersionTag] = strings.ToLower(m.Version)
		inferredTags = append(inferredTags, langVersionTag)
	}

	(*annotations)["effx.io/inferred-tags"] = strings.Join(inferredTags, ",")

	config.SetTags(*tags)
	config.SetAnnotations(*annotations)

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

	relativePath, err := convertToRelativePath(y.FilePath)
	if err != nil {
		relativePath = y.FilePath
	}

	versionControlUrl := getVersionControlLink(y.FilePath, relativePath)

	config.FileContents = string(yamlFile)
	config.SetAnnotations(map[string]string{
		"effx.io/source":               "effx-cli",
		"effx.io/file-path":            relativePath,
		"effx.io/version-control-link": versionControlUrl,
		"effx.io/repository-name":      getRepoName(versionControlUrl),
	})

	if os.Getenv("DISABLE_LANGUAGE_DETECTION") != "true" {
		result, err := metadata.InferMetadata(filepath.Dir(y.FilePath))
		if err != nil {
			log.Printf("Could not predict version %+v\n", err)
		}
		config = setMetadata(config, result)
	}

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

	url := GenerateUrl()
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

	url := GenerateUrl()
	url.Path = "v2/config"

	request, _ := http.NewRequest("PUT", url.String(), bytes.NewReader(body))
	request.Header.Add("content-type", "application/json")
	request.Header.Add("x-effx-api-key", apiKey)

	// resp, err := http.DefaultClient.Do(request)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

	// return checkForErrors(resp)
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GenerateUrl() *url.URL {
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
