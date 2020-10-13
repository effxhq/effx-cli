package parser

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/effxhq/effx-cli/data"
)

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

func ProcessArgs(filePath string, directory string) ([]data.ApiResource, error) {
	var resources []data.ApiResource
	var err error
	if filePath != "" {
		resources, err = ProcessFile(filePath)
	} else {
		resources, err = ProcessDirectory(directory)
	}
	if err != nil {
		return nil, err
	}
	return resources, err
}

func ProcessFile(filePath string) ([]data.ApiResource, error) {
	effxYaml := EffxYaml{FilePath: filePath}
	return effxYaml.ToApiResources()
}

func ProcessDirectory(directory string) ([]data.ApiResource, error) {
	pattern := EffxYaml{}.getFilePattern()
	matches, _ := glob(directory, pattern)

	var resources []data.ApiResource
	for _, path := range matches {
		fileResources, err := ProcessFile(path)
		resources = append(resources, fileResources...)

		if err != nil {
			return nil, err
		}
	}

	return resources, nil
}
