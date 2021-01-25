package parser

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	effx_v1_api "github.com/effxhq/effx-api/generated/go"
	"github.com/effxhq/effx-cli/data"
)

type EventPayload struct {
	Name            string
	Description     string
	ServiceName     string
	IntegrationName string
	ImageUrl        string
	Email           string
	Tags            string
	Hashtags        string
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

func ProcessEvent(e *EventPayload) data.V1Data {
	tagsPayload := []effx_v1_api.TagPayload{}
	hashtagsPayload := []string{}

	if e.Tags != "" {
		tagsStringNoSpace := strings.Join(strings.Fields(e.Tags), "")
		splitTagsString := strings.Split(tagsStringNoSpace, ",")

		for _, splitTag := range splitTagsString {
			splitTagString := strings.Split(splitTag, ":")

			if len(splitTagString) == 2 {
				tagsPayload = append(tagsPayload, effx_v1_api.TagPayload{Key: splitTagString[0], Value: splitTagString[1]})
			} else {
				log.Fatalf("found invalid tag: %s", splitTag)
			}
		}
	}

	if e.Hashtags != "" {
		hashtagsStringNoSpace := strings.Join(strings.Fields(e.Hashtags), "")
		hashtagsPayload = strings.Split(hashtagsStringNoSpace, ",")
	}

	object := &data.V1Data{
		Event: &effx_v1_api.EventPayload{
			ProducedAtTimeMilliseconds: time.Now().UnixNano() / 1e6,
			Name:                       e.Name,
			Description:                e.Description,
			Tags:                       tagsPayload,
			Hashtags:                   hashtagsPayload,
			Integration: &effx_v1_api.IntegrationPayload{
				Name: e.IntegrationName,
			},
		},
	}

	if e.Email != "" {
		object.Event.User = &effx_v1_api.EventUserPayload{
			Email: e.Email,
		}
	}
	return *object
}
