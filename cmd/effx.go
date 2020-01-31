package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/effxhq/effx-go/cmd/sync"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "effx",
	Short: "effx cli client",
	Long:  `effx cli client`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello world")
	},
}

func Initialize() {
	sync.Initialize()
	rootCmd.AddCommand(sync.SyncCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
