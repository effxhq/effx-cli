package metadata

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/effxhq/effx-cli/metadata/golang"
	"github.com/effxhq/effx-cli/metadata/java"
	"github.com/effxhq/effx-cli/metadata/javascript"
	"github.com/effxhq/effx-cli/metadata/php"
)

type Result struct {
	Language string
	Version  string
}

type FileHandler struct {
	Langugage string
	FileName  string
	Handler   func(fileContent string) string
}

func (fh *FileHandler) GetResult(fileContent string) *Result {
	return &Result{
		Language: fh.Langugage,
		Version:  fh.Handler(fileContent),
	}
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

// there could be multiple files that can
// contain relevant information
func getRelevantFiles(lang string) []*FileHandler {
	switch lang {
	case "go":
		return []*FileHandler{
			{
				Langugage: "go",
				FileName:  "go.mod",
				Handler:   golang.HandleGoModFile,
			}}
	case "javascript":
		return []*FileHandler{
			{
				Langugage: "node",
				FileName:  "package.json",
				Handler:   javascript.HandlePackageJson,
			}}
	case "php":
		return []*FileHandler{
			{
				Langugage: "php",
				FileName:  "composer.json",
				Handler:   php.HandleComposerJson,
			}}
	case "java":
		return []*FileHandler{
			{
				Langugage: "java",
				FileName:  "pom.xml",
				Handler:   java.HandlePomFile,
			}}
	default:
		return []*FileHandler{}
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

				for _, fileHandler := range relavantFiles {
					if file.Name() == fileHandler.FileName {
						res = fileHandler.GetResult(string(content))
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
