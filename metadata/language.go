package metadata

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-enry/go-enry/v2"
)

// EffxYamlPattern is the regex pattern for yaml files
const EffxYamlPattern = "(.+\\.)?effx\\.ya?ml$"

var effxYamlRegex = regexp.MustCompile(EffxYamlPattern)

const NonCodePattern = `.yaml|.yml|.json|.md`

var fileIgnorePattern = regexp.MustCompile(NonCodePattern)

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

// inferLanguage detects the programming language used in the provided directory.
// this function is passed the directories of where effx.yaml are found
func inferLanguage(pathDir string) (string, error) {
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

			// other files to ignore
			if fileIgnorePattern.MatchString(fileName) {
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

	err := filepath.Walk(pathDir, collector)

	return determineMostCommonLangugage(languageCount), err
}
