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

var (
	apiKeyString    string
	directoryString string
	filePathString  string
)

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync effx.yaml to effx api",
	Long:  `sync effx.yaml to effx api`,
	Run: func(cmd *cobra.Command, args []string) {
		if apiKeyString == "" {
			log.Fatal("api key is required")
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

func Initialize() {
	SyncCmd.PersistentFlags().StringVarP(&apiKeyString, "api_key", "k", "", "api_key")
	SyncCmd.PersistentFlags().StringVarP(&filePathString, "file_path", "f", "", "file_path")
	SyncCmd.PersistentFlags().StringVarP(&directoryString, "directory", "d", "", "directory")
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
		object, err := parser.YamlFile(filePath)

		if err != nil {
			return err
		}

		if err := validator.ValidateObject(object); err != nil {
			return err
		}

		log.Printf("sending %s to effx api", filePath)

		return c.Synchronize(object)
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
