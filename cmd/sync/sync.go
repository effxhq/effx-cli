package sync

import (
	"errors"
	"log"
	"os"

	"github.com/effxhq/effx-cli/internal/parser"
	"github.com/spf13/cobra"
)

const effxApiKeyName = "EFFX_API_KEY"

var (
	apiKeyString             string
	directoryString          string
	filePathString           string
	disableLanguageDetection bool
)

func init() {
	SyncCmd.PersistentFlags().StringVarP(&apiKeyString, "key", "k", "", "your effx api key. alternatively, you can use env var EFFX_API_KEY")
	SyncCmd.PersistentFlags().StringVarP(&filePathString, "file", "f", "", "path to a effx.yaml file")
	SyncCmd.PersistentFlags().StringVarP(&directoryString, "dir", "d", "", "directory to recursively find and sync effx.yaml files")
	SyncCmd.PersistentFlags().BoolVarP(&disableLanguageDetection, "disable-languauge-detection", "", false, "disables automatic langugage and version detections")
}

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync effx.yaml file(s) to the effx api",
	Long:  `sync effx.yaml file(s) to the effx api`,
	Args: func(cmd *cobra.Command, args []string) error {
		if apiKeyString == "" {
			if apiKeyString = os.Getenv(effxApiKeyName); apiKeyString == "" {
				return errors.New("api key is required")
			}
		}

		if filePathString == "" && directoryString == "" {
			return errors.New("-f <file_path> or -d <directory> is required")
		}

		if filePathString != "" && directoryString != "" {
			return errors.New("-f <file_path> and -d <directory> cannot be used together")
		}

		if disableLanguageDetection {
			os.Setenv("DISABLE_LANGUAGE_DETECTION", "true")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		resources := parser.ProcessArgs(filePathString, directoryString)

		filePaths := []string{}

		for _, file := range resources {
			filePaths = append(filePaths, file.FilePath)
		}

		err := parser.DetectServicesFromEffxYamls(filePaths, apiKeyString, "effx-cli")
		if err != nil {
			log.Println("Could not send detected services, err:", err)
		}

		for _, resource := range resources {
			err := resource.Sync(apiKeyString)
			if err != nil {
				return err
			}
		}

		return nil
	},
}
