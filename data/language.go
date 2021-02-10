package data

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-enry/go-enry/v2"
)

func determineMostCommonLangugage(languageCount map[string]int) string {
	max := 0
	mostCommonLang := ""

	for key, value := range languageCount {
		if max < value {
			max = value
			mostCommonLang = key
		}
	}

	return mostCommonLang
}

// InferLanguage detects the programming used in the provided work directory .
func inferLanguage(workDir string) (string, error) {
	languageCount := map[string]int{}

	collector := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if !info.IsDir() {
			fileName := filepath.Base(path)

			// we don't want to look at effx files
			if effxYamlRegex.MatchString(fileName) {
				return nil
			}

			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// infers langugage from extension, and code content.
			lang := strings.ToLower(enry.GetLanguage(fileName, content))

			if count, ok := languageCount[lang]; ok {
				languageCount[lang] = count + 1
			} else {
				languageCount[lang] = 1
			}
		}
		return nil
	}

	err := filepath.Walk(workDir, collector)

	return determineMostCommonLangugage(languageCount), err
}
