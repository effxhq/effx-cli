package cmd

import (
	"log"
	"os"

	"github.com/effxhq/effx-cli/cmd/event"
	"github.com/effxhq/effx-cli/cmd/lint"
	"github.com/effxhq/effx-cli/cmd/sync"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "effx",
	Short: "effx cli client",
	Long:  `effx cli client`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("this is the effx cli")
	},
}

func init() {
	rootCmd.AddCommand(lint.LintCmd)
	rootCmd.AddCommand(sync.SyncCmd)
	rootCmd.AddCommand(event.EventCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
