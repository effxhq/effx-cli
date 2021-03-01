package parser

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	effx_api "github.com/effxhq/effx-api-v2/generated/go/client"
	"github.com/effxhq/effx-cli/data"
	"github.com/effxhq/effx-cli/internal/discover"
)

type EventPayload struct {
	Title            string
	Message          string
	ServiceName      string
	Tags             []string
	Actions          []string
	ProducedAtTimeMS int
}

func glob(dir string, pattern string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		matched, _ := regexp.MatchString(pattern, path)
		if matched {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func ProcessArgs(filePath string, directory string) []data.EffxYaml {
	var yamls []data.EffxYaml
	if filePath != "" {
		yamls = ProcessFile(filePath)
	} else {
		yamls = ProcessDirectory(directory)
	}
	return yamls
}

func ProcessFile(filePath string) []data.EffxYaml {
	effxYaml := data.EffxYaml{FilePath: filePath}
	return []data.EffxYaml{effxYaml}
}

func ProcessDirectory(directory string) []data.EffxYaml {
	matches, _ := glob(directory, data.EffxYamlPattern)

	var yamls []data.EffxYaml
	for _, path := range matches {
		yamls = append(yamls, data.EffxYaml{FilePath: path})
	}

	return yamls
}

func DetectServicesFromEffxYamls(filePaths []string, apiKeyString, sourceName string) error {
	services := discover.DetectServices(sourceName, filePaths)

	return discover.SendDetectedServices(apiKeyString, data.GenerateUrl(), services)
}

func ProcessEvent(e *EventPayload) *data.EffxEvent {
	tagsPayload := []effx_api.CreateEventPayloadTags{}
	actions := []effx_api.CreateEventPayloadActions{}
	producedAtTime := int64(e.ProducedAtTimeMS)

	if len(e.Tags) > 0 {
		for _, tag := range e.Tags {
			splitTagString := strings.Split(tag, ":")

			if len(splitTagString) == 2 {
				tagsPayload = append(tagsPayload, effx_api.CreateEventPayloadTags{
					Key:   splitTagString[0],
					Value: splitTagString[1],
				})
			} else {
				log.Fatalf("found invalid tag: %s", tag)
			}
		}
	}

	if len(e.Actions) > 0 {
		for _, action := range e.Actions {
			// format: level:name:url
			res := strings.SplitN(action, ":", 3)

			if len(res) < 2 {
				log.Fatalf("found invalid action: %s", action)
			}

			actions = append(actions, effx_api.CreateEventPayloadActions{
				Level: res[0],
				Name:  res[1],
				Url:   res[2],
			})
		}
	}

	payload := &data.EffxEvent{
		Payload: &effx_api.CreateEventPayload{
			Title:       e.Title,
			Message:     e.Message,
			ServiceName: &e.ServiceName,
			Tags:        &tagsPayload,
			Actions:     &actions,
		},
	}

	// if optional produced at timstamp is less than a year ago.
	if producedAtTime > time.Now().AddDate(-1, 0, 0).UnixNano()/1e6 {
		payload.Payload.TimestampMilliseconds = &producedAtTime
	}

	return payload
}
