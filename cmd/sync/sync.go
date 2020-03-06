package sync

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/effxhq/effx-go/internal/client"
	"github.com/effxhq/effx-go/internal/client/http"
	"github.com/effxhq/effx-go/internal/parser"
	"github.com/effxhq/effx-go/internal/validator"
	"github.com/spf13/cobra"
)

const effxApiKeyName = "EFFX_API_KEY"

var (
	apiKeyString    string
	directoryString string
	filePathString  string
	isDryRun        bool
)

func Initialize() {
	SyncCmd.PersistentFlags().StringVarP(&apiKeyString, "key", "k", "", "your effx api key. alternatively, you can use env var EFFX_API_KEY")
	SyncCmd.PersistentFlags().StringVarP(&filePathString, "file", "f", "", "path to a effx.yaml file")
	SyncCmd.PersistentFlags().StringVarP(&directoryString, "dir", "d", "", "directory to recursively find and sync effx.yaml files")
	SyncCmd.PersistentFlags().BoolVarP(&isDryRun, "dry-run", "", false, "validate file(s)")
}

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync effx.yaml file(s) to the effx api",
	Long:  `sync effx.yaml file(s) to the effx api`,
	Run: func(cmd *cobra.Command, args []string) {
		if apiKeyString == "" && !isDryRun {
			if apiKeyString = os.Getenv(effxApiKeyName); apiKeyString == "" {
				log.Fatal("api key is required")
			}
		}

		if filePathString == "" && directoryString == "" {
			log.Fatal("-f <file_path> or -d <directory> is required")
		}

		if filePathString != "" && directoryString != "" {
			log.Fatal("-f <file_path> and -d <directory> cannot be used together")
		}

		c := http.New(apiKeyString)

		if filePathString != "" {
			if matched, err := isEffxYaml(filePathString); err != nil {
				log.Fatalf("unexpected error: %v", err)
			} else {
				if matched {
					if err := processFile(filePathString, c); err != nil {
						log.Fatalf("error: %v", err)
					}
				} else {
					log.Fatalf("file_path %s is invalid", filePathString)
				}
			}
		} else {
			if err := processDirectory(directoryString, c); err != nil {
				log.Fatalf("error: %v", err)
			}

		}
	},
}

func isEffxYaml(filePath string) (bool, error) {
	matched, err := regexp.MatchString(`.*effx\.yaml$`, filePath)

	return matched, err
}

func processFile(filePath string, c client.Client) error {
	ok, err := isEffxYaml(filePath)

	if err != nil {
		return err
	}

	if ok {
		log.Printf("parsing %s", filePath)
		object, err := parser.YamlFile(filePath)

		if err != nil {
			return err
		}

		if err := validator.ValidateObject(object); err != nil {
			return err
		}

		if isDryRun {
			log.Printf("%s is valid", filePath)
		} else {
			log.Printf("sending %s to effx api", filePath)

			return c.Synchronize(object)
		}
	}

	return nil
}

func processDirectory(directory string, c client.Client) error {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		err = processFile(path, c)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
