package metadata

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/effxhq/effx-cli/metadata/golang"
	"github.com/effxhq/effx-cli/metadata/javascript"
)

type Result struct {
	Language string
	Version  string
}

func New(lang, version string) *Result {
	return &Result{
		Language: lang,
		Version:  version,
	}
}

func isRootDirectory(files []os.FileInfo) bool {
	for _, file := range files {
		if file.Name() == ".git" {
			return true
		}
	}
	return false
}

func getRelevantFiles(lang string) []string {
	switch lang {
	case "go":
		return []string{"go.mod"}
	case "javascript":
		return []string{"package.json"}
	default:
		return []string{}
	}
}

func getVersion(lang string, fileContent string) *Result {
	switch lang {
	case "go":
		return New("go", golang.HandleGoModFile(fileContent))
	case "javascript":
		return New("node", javascript.HandlePackageJson(fileContent))
	default:
		return nil
	}
}

func InferMetadata(pathDir string) (*Result, error) {
	lang, err := inferLanguage(pathDir)
	if err != nil {
		return nil, err
	}

	version, err := inferVersion(pathDir, lang)
	if err != nil {
		return nil, err

	}

	return version, nil
}

func inferVersion(pathDir string, lang string) (*Result, error) {
	var res *Result

	for pathDir != "" {
		files, err := ioutil.ReadDir(pathDir)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if !file.IsDir() {
				content, err := ioutil.ReadFile(pathDir + "/" + file.Name())
				if err != nil {
					return nil, err
				}

				relavantFiles := getRelevantFiles(lang)

				for _, relavantFile := range relavantFiles {
					if file.Name() == relavantFile {
						res = getVersion(lang, string(content))
						if res != nil {
							return res, nil
						}
					}
				}
			}
		}

		if isRootDirectory(files) {
			return nil, nil
		}

		pathDir = filepath.Join(pathDir, "..")
	}

	return nil, nil
}
