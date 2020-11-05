package lint

import (
	"errors"

	"github.com/effxhq/effx-cli/internal/parser"
	"github.com/spf13/cobra"
)

var (
	directoryString string
	filePathString  string
)

func init() {
	LintCmd.PersistentFlags().StringVarP(&filePathString, "file", "f", "", "path to a effx.yaml file")
	LintCmd.PersistentFlags().StringVarP(&directoryString, "dir", "d", "", "directory to recursively find and sync effx.yaml files")
}

var LintCmd = &cobra.Command{
	Use:   "lint",
	Short: "lint effx.yaml file(s) to the effx api",
	Long:  `lint effx.yaml file(s) to the effx api`,
	Args: func(cmd *cobra.Command, args []string) error {
		if filePathString == "" && directoryString == "" {
			return errors.New("-f <file_path> or -d <directory> is required")
		}

		if filePathString != "" && directoryString != "" {
			return errors.New("-f <file_path> and -d <directory> cannot be used together")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		resources := parser.ProcessArgs(filePathString, directoryString)
		for _, resource := range resources {
			err := resource.Lint()
			if err != nil {
				return err
			}
		}
		return nil
	},
}
