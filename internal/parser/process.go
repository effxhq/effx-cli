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
)

type EventPayload struct {
	Title            string
	Name             string
	Message          string
	ServiceName      string
	Tags             string
	Actions          string
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

func ProcessEvent(e *EventPayload) *data.EffxEvent {
	tagsPayload := []effx_api.CreateEventPayloadTags{}
	actions := []effx_api.CreateEventPayloadActions{}
	timestampMilliseconds := time.Now().UnixNano() / 1e6
	producedAtTime := int64(e.ProducedAtTimeMS)

	if e.Tags != "" {
		tagsStringNoSpace := strings.Join(strings.Fields(e.Tags), "")
		splitTagsString := strings.Split(tagsStringNoSpace, ",")

		for _, splitTag := range splitTagsString {
			splitTagString := strings.Split(splitTag, ":")

			if len(splitTagString) == 2 {
				tagsPayload = append(tagsPayload, effx_api.CreateEventPayloadTags{
					Key:   splitTagString[0],
					Value: splitTagString[1],
				})
			} else {
				log.Fatalf("found invalid tag: %s", splitTag)
			}
		}
	}

	if e.Actions != "" {
		// format: level:name:url
		res := strings.SplitN(e.Actions, ":", 3)

		if len(res) < 2 {
			log.Fatalf("found invalid tag: %s", e.Actions)
		}

		actions = append(actions, effx_api.CreateEventPayloadActions{
			Level: res[0],
			Name:  res[1],
			Url:   res[2],
		})
	}

	// if optional produced at timstamp is less than a year ago.
	if producedAtTime < time.Now().AddDate(-1, 0, 0).UnixNano()/1e6 {
		timestampMilliseconds = producedAtTime
	}

	payload := &data.EffxEvent{
		Payload: &effx_api.CreateEventPayload{
			Title:                 e.Title,
			Message:               e.Message,
			Tags:                  &tagsPayload,
			Actions:               &actions,
			TimestampMilliseconds: timestampMilliseconds,
		},
	}

	return payload
}
