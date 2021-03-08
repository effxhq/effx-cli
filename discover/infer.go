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

var serviceDirectoryNames = []string{
	"services",
	"apps",
}

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

// DetectServicesFromRelavantFiles detects services based on
// containing a service-like file (package.json etc)
func DetectServicesFromRelavantFiles(workdir string, effxFiles []data.EffxYaml, sourceName string) ([]effx_api.DetectedServicesPayload, error) {
	detectedServices := []effx_api.DetectedServicesPayload{}
	effxFileLocations := filePathsFromEffxYaml(effxFiles)

	err := filepath.Walk(workdir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			dirName := f.Name()

			if directoryContainsEffxYaml(f, effxFileLocations) {
				return nil
			}

			// if directory name contain "services", "apps" etc.
			for _, relavantName := range serviceDirectoryNames {
				if strings.Contains(dirName, relavantName) {
					detectedServices = append(detectedServices, createDetectedServicePayload(f, sourceName, path))
					return nil
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
