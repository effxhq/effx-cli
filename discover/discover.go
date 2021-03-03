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

	effx_api "github.com/effxhq/effx-api-v2/generated/go/client"
	"github.com/effxhq/effx-cli/data"
	"github.com/effxhq/effx-cli/internal/parser"
	"github.com/effxhq/effx-cli/metadata"
)

func DetectServicesFromWorkDir(workDir string, apiKeyString, sourceName string) error {
	filePaths := parser.ProcessDirectory(workDir)
	return DetectServicesFromEffxYamls(filePaths, apiKeyString, sourceName)
}

func DetectServicesFromEffxYamls(files []data.EffxYaml, apiKeyString, sourceName string) error {
	filePaths := []string{}

	for _, file := range files {
		filePaths = append(filePaths, file.FilePath)
	}

	services := DetectServices(sourceName, filePaths)
	return SendDetectedServices(apiKeyString, data.GenerateUrl(), services)
}

func findCommonDirectory(effxFileLocations []string) string {
	matchedEffxFiles := generateIterators(effxFileLocations)
	prefixString := ""

	for len(matchedEffxFiles) > 0 {
		count := make(map[string]int)
		for _, matchedFile := range matchedEffxFiles {
			peek := matchedFile.Peek()
			if peek != "" {
				count[peek]++
			}
		}

		maxK := ""
		maxV := 1
		for k, v := range count {
			if v > maxV {
				maxK = k
				maxV = v
			}
		}

		nextRound := make([]*Iterator, 0, maxV)
		for _, matchedFile := range matchedEffxFiles {
			// advance ptr
			if matchedFile.Next() == maxK {
				// put into next
				nextRound = append(nextRound, matchedFile)
			}
		}

		prefixString += maxK
		matchedEffxFiles = nextRound

	}

	if prefixString == "" {
		return ""
	}

	// prefix string should be a directory ending with a slash
	slashIndex := strings.LastIndex(prefixString, "/")

	if slashIndex != len(prefixString) {
		// trim file name, keep last dir slash
		// example:
		// services/dooku -> services/
		prefixString = prefixString[:slashIndex+1]
	}

	return prefixString
}

// returns a list of names of detected services
func DetectServices(sourceName string, effxFileLocations []string) []effx_api.DetectedServicesPayload {
	detectedServices := []effx_api.DetectedServicesPayload{}

	commonDir := findCommonDirectory(effxFileLocations)

	files, err := ioutil.ReadDir(commonDir)
	if err != nil {
		return detectedServices
	}

	for _, file := range files {
		// looking at directories only for service locations
		if file.IsDir() {
			contains := false
			for _, effxFileLocation := range effxFileLocations {
				if strings.Contains(effxFileLocation, file.Name()) {
					contains = true
				}
			}
			if !contains {
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
