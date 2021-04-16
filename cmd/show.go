package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xunull/goc/commonx"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show",
	Long:  `show`,
}

var showWordsCmd = &cobra.Command{
	Use:   "words",
	Short: "words",
	Long:  `words`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		res := commonx.GetWords(target)
		for k, _ := range res {
			fmt.Println(k)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.AddCommand(showWordsCmd)
}
