package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type PackageJson struct {
	Engines struct {
		Node string `json:"node"`
	} `json:"engines"`
}

type result struct {
	lang    string
	version string
}

// regex to find go version from go.mod file
const GoVersionRegex = `go \d+(\.\d+)+`

var goModRegex = regexp.MustCompile(GoVersionRegex)

func isRootDirectory(files []os.FileInfo) bool {
	for _, file := range files {
		if file.Name() == ".git" {
			return true
		}
	}
	return false
}

func handleGoModFile(fileContent string) *result {
	res := goModRegex.FindString(fileContent)
	versionStr := strings.Split(res, " ")

	if len(versionStr) < 2 {
		return nil
	}

	return &result{
		lang:    "go",
		version: versionStr[1],
	}
}

func handlePackageJson(fileContent string) *result {
	var packageJson = &PackageJson{}

	err := json.Unmarshal([]byte(fileContent), packageJson)
	if err != nil {
		return nil
	}

	return &result{
		lang:    "node",
		version: packageJson.Engines.Node,
	}
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

func getVersion(lang string, fileContent string) *result {
	switch lang {
	case "go":
		return handleGoModFile(fileContent)
	case "javascript":
		return handlePackageJson(fileContent)
	default:
		return nil
	}
}

func inferMetadata(pathDir string) (*result, error) {
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

func inferVersion(pathDir string, lang string) (*result, error) {
	var res *result

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
