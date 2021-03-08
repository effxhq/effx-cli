package discover

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	effx_api "github.com/effxhq/effx-api-v2-go/client"
	"github.com/effxhq/effx-cli/data"
	"github.com/effxhq/effx-cli/internal/parser"
	"github.com/effxhq/effx-cli/metadata"
)

func filePathsFromEffxYaml(files []data.EffxYaml) []string {
	filePaths := []string{}

	for _, file := range files {
		filePaths = append(filePaths, file.FilePath)
	}
	return filePaths
}

func directoryContainsEffxYaml(file os.FileInfo, effxFileLocations []string) bool {
	if !file.IsDir() {
		return false
	}

	contains := false
	for _, effxFileLocation := range effxFileLocations {
		if strings.Contains(effxFileLocation, file.Name()) {
			contains = true
		}
	}
	return contains
}

func createDetectedServicePayload(file os.FileInfo, sourceName string, commonDir string) effx_api.DetectedServicesPayload {
	payload := effx_api.DetectedServicesPayload{
		Name:       strings.ToLower(file.Name()),
		SourceName: &sourceName,
		Tags:       &map[string]string{},
	}

	if os.Getenv("DISABLE_LANGUAGE_DETECTION") != "true" {
		result, err := metadata.InferMetadata(filepath.Dir(commonDir + file.Name()))
		if err != nil {
			log.Printf("Could not predict version %+v\n", err)
		}

		if result != nil {
			if result.Language != "" {
				language := strings.ToLower(result.Language)
				(*payload.Tags)["language"] = language

				// version cannot be found if language cannot be detected
				if result.Version != "" {
					(*payload.Tags)[language] = strings.ToLower(result.Version)
				}
			}
		}
	}

	return payload
}

func DetectServicesFromWorkDir(workDir string, apiKeyString, sourceName string) error {
	filePaths := parser.ProcessDirectory(workDir)
	services, err := DetectServicesFromRelavantFiles(workDir, filePaths, sourceName)
	if err != nil {
		return err
	}

	servicesInferredFromYaml := DetectServicesFromEffxYamls(filePaths, apiKeyString, sourceName)
	services = append(services, servicesInferredFromYaml...)

	return SendDetectedServices(apiKeyString, data.GenerateUrl(), services)
}

// returns a list of names of detected services

func DetectServicesFromEffxYamls(effxFiles []data.EffxYaml, apiKeyString, sourceName string) []effx_api.DetectedServicesPayload {
	effxFileLocations := filePathsFromEffxYaml(effxFiles)
	detectedServices := []effx_api.DetectedServicesPayload{}

	commonDir := findCommonDirectory(effxFileLocations)

	files, err := ioutil.ReadDir(commonDir)
	if err != nil {
		return detectedServices
	}

	for _, file := range files {
		// looking at directories only for service locations
		if file.IsDir() {
			contains := directoryContainsEffxYaml(file, effxFileLocations)

			if !contains {
				payload := createDetectedServicePayload(file, sourceName, commonDir)

				detectedServices = append(detectedServices, payload)
			}
		}
	}

	return detectedServices
}

func SendDetectedServices(apiKey string, url *url.URL, servicePayloads []effx_api.DetectedServicesPayload) error {
	for _, payload := range servicePayloads {
		body, _ := json.Marshal(payload)

		url.Path = "v2/detected_services"

		request, _ := http.NewRequest("PUT", url.String(), bytes.NewReader(body))
		request.Header.Add("content-type", "application/json")
		request.Header.Add("x-effx-api-key", apiKey)

		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}

	log.Println("Successfully detected ", len(servicePayloads), " service(s)")
	return nil
}
