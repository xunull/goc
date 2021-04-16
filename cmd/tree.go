package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xunull/goc/commonx"
	"github.com/xunull/goc/tree"
	"os"
)

// todo error
// treeCmd represents the tree command
var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "like tree",
	Long:  `like tree`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := os.Getwd()
		commonx.CheckErrOrFatal(err)
		ppt := tree.DirTree(p)
		res := tree.TreeIt([]tree.TreeAble{ppt})
		for _, item := range res {
			color.Blue(item)
		}
	},
}

func init() {
	rootCmd.AddCommand(treeCmd)

}
