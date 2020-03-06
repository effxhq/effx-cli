package cmd

import (
	"log"
	"os"

	"github.com/effxhq/effx-go/cmd/event"
	"github.com/effxhq/effx-go/cmd/sync"
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

func Initialize() {
	sync.Initialize()
	event.Initialize()

	rootCmd.AddCommand(sync.SyncCmd)
	rootCmd.AddCommand(event.EventCreateCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
