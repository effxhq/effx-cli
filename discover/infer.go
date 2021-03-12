package discover

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	effx_api "github.com/effxhq/effx-api-v2-go/client"
	"github.com/effxhq/effx-cli/data"
	"github.com/thoas/go-funk"
)

var relavantFiles = []string{
	"go.mod",           // golang/gin
	"package.json",     // nodejs
	"requirements.txt", // python/django/flask
	"pom.xml",          // java/maven
	"gemfile",          // rails, sinatra
	"composer.json",    // larvel
	"build.gradle",     // java spring
	"mix.exs",          // elixer/phoenix
}

var defaultInferredServiceDirectoryNames = []string{
	"services",
	"apps",
}

func getInferredServiceDirectoryNames() []string {
	configuredDirNames := os.Getenv("INFERRED_SERVICE_DIRECTORY_NAMES")

	if configuredDirNames != "" {
		return strings.Split(configuredDirNames, ",")
	}

	return defaultInferredServiceDirectoryNames
}

// DetectServicesFromFiles detects services based on
// containing a service-like file (package.json etc)
func DetectServicesFromFiles(workdir string, effxFiles []data.EffxYaml, sourceName string) ([]effx_api.DetectedServicesPayload, error) {
	detectedServices := []effx_api.DetectedServicesPayload{}
	effxFileLocations := filePathsFromEffxYaml(effxFiles)

	err := filepath.Walk(workdir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			dirName := f.Name()

			if directoryContainsEffxYaml(f, effxFileLocations) {
				return nil
			}

			// if directory name contains for example: "services", "apps" etc,
			// all subdirectories are services.
			for _, relavantName := range getInferredServiceDirectoryNames() {
				if strings.Contains(dirName, relavantName) {
					files, err := ioutil.ReadDir(path + "/" + f.Name())
					if err == nil {
						for _, file := range files {
							if file.IsDir() {
								detectedServices = append(detectedServices, createDetectedServicePayload(f, sourceName, path))
							}
						}
						return nil
					}
				}
			}

			files, err := ioutil.ReadDir(path)
			if err != nil {
				return err
			}

			for _, file := range files {
				// if directory contains a relavant service file
				if funk.ContainsString(relavantFiles, file.Name()) {
					detectedServices = append(detectedServices, createDetectedServicePayload(f, sourceName, path))
					return nil
				}
			}
		}
		return nil
	})

	return detectedServices, err
}
