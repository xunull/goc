package cmd

import (
	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "random",
	Long:  `random`,
}

var randomMakeFilesCmd = &cobra.Command{
	Use:   "makefiles",
	Short: "make files",
	Long:  `make files`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//width, _ := cmd.Flags().GetInt("width")
		//depth, _ := cmd.Flags().GetInt("depth")
		//dirname := args[0]
		////file_utils.MakeTempFiles(dirname, width, depth)
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)
	randomCmd.AddCommand(randomMakeFilesCmd)

	randomMakeFilesCmd.Flags().Int("width", 10, "file count")
	randomMakeFilesCmd.Flags().Int("depth", 1, "dir depth")
}
