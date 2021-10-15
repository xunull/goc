package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// wordCmd represents the word command
var wordCmd = &cobra.Command{
	Use:   "word",
	Short: "word",
	Long:  `word`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("word called")
	},
}

func init() {
	rootCmd.AddCommand(wordCmd)
}
