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
