package parser

import (
	"path/filepath"

	"github.com/effxhq/effx-cli/data"
)

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
	pattern := "**/*.effx.yaml"
	yaml, _ := filepath.Glob(pattern)

	pattern = "**/*.effx.yml"
	yml, _ := filepath.Glob(pattern)

	matches := append(yml, yaml...)

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
