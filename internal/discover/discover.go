package discover

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	effx_api "github.com/effxhq/effx-api-v2/generated/go/client"
)

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
func DetectServices(effxFileLocations []string) []string {
	detectedServiceNames := []string{}

	commonDir := findCommonDirectory(effxFileLocations)

	files, err := ioutil.ReadDir(commonDir)
	if err != nil {
		return []string{}
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
				detectedServiceNames = append(detectedServiceNames, file.Name())
			}
		}
	}

	return detectedServiceNames
}

func SendDetectedServices(apiKey, sourceName string, url *url.URL, services []string) error {

	for _, serviceName := range services {
		payload := effx_api.DetectedServicesPayload{
			Name:       serviceName,
			SourceName: &sourceName,
		}
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

	log.Println("Successfully detected ", len(services), " service(s)")
	return nil
}
