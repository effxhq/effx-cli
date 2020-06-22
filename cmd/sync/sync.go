package sync

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/antihax/optional"
	effx_api "github.com/effxhq/effx-api/generated/go"
	"github.com/effxhq/effx-cli/internal/parser"
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
		if apiKeyString == "" {
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

		if filePathString != "" {
			if matched, err := isEffxYaml(filePathString); err != nil {
				log.Fatalf("unexpected error: %v", err)
			} else {
				if matched {
					if err := processFile(filePathString); err != nil {
						log.Fatalf("error: %v", err)
					}
				} else {
					log.Fatalf("file_path %s is invalid", filePathString)
				}
			}
		} else {
			if err := processDirectory(directoryString); err != nil {
				log.Fatalf("error: %v", err)
			}

		}
	},
}

func isEffxYaml(filePath string) (bool, error) {
	matched, err := regexp.MatchString(`.*effx\.yaml$`, filePath)

	return matched, err
}

func processFile(filePath string) error {
	ok, err := isEffxYaml(filePath)

	if err != nil {
		return err
	}

	if ok {
		log.Printf("parsing %s", filePath)
		objects, err := parser.YamlFile(filePath)

		if err != nil {
			return err
		}

		log.Printf("sending %s to effx api", filePath)

		basePath := "https://api.effx.io/v1"

		if os.Getenv("EFFX_API_HOST") != "" {
			log.Printf("switching to use basePath %s", fmt.Sprintf("%s/v1", os.Getenv("EFFX_API_HOST")))
			basePath = fmt.Sprintf("%s/v1", os.Getenv("EFFX_API_HOST"))
		}

		client := effx_api.NewAPIClient(&effx_api.Configuration{
			BasePath:      basePath,
			DefaultHeader: make(map[string]string),
			UserAgent:     "Swagger-Codegen/1.0.0/go",
		})

		for _, obj := range objects {
			if obj.Service != nil {
				_, err := client.ServicesApi.ServicesPut(
					context.Background(),
					apiKeyString,
					*obj.Service,
					&effx_api.ServicesPutOpts{
						XEffxValidateOnly: optional.NewString(fmt.Sprintf("%v", isDryRun)),
					},
				)

				if err != nil {
					log.Printf("error %s", string(err.(effx_api.GenericSwaggerError).Body()))

					return err
				}
			} else if obj.User != nil {
				_, err := client.UsersApi.UsersPut(
					context.Background(),
					apiKeyString,
					*obj.User,
					&effx_api.UsersPutOpts{
						XEffxValidateOnly: optional.NewString(fmt.Sprintf("%v", isDryRun)),
					},
				)

				if err != nil {
					log.Printf("error %s", string(err.(effx_api.GenericSwaggerError).Body()))

					return err
				}
			} else if obj.Team != nil {
				_, err := client.TeamsApi.TeamsPut(
					context.Background(),
					apiKeyString,
					*obj.Team,
					&effx_api.TeamsPutOpts{
						XEffxValidateOnly: optional.NewString(fmt.Sprintf("%v", isDryRun)),
					},
				)

				if err != nil {
					log.Printf("error %s", string(err.(effx_api.GenericSwaggerError).Body()))

					return err
				}
			} else {
				return fmt.Errorf("unsupported object type found, %v", obj)
			}
		}
	}

	return nil
}

func processDirectory(directory string) error {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		err = processFile(path)

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
