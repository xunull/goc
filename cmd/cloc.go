package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// clocCmd represents the cloc command
var clocCmd = &cobra.Command{
	Use:   "cloc",
	Short: "like cloc",
	Long:  `like cloc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cloc called")
	},
}

func init() {
	rootCmd.AddCommand(clocCmd)
}
