package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xunull/goc/commonx"
	"os"
	"path/filepath"
)

// md5Cmd represents the md5 command
var md5Cmd = &cobra.Command{
	Use:   "md5",
	Short: "md5",
	Long:  `md5`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		t := args[0]

		if filepath.IsAbs(t) {
			res, err := commonx.MD5File(t)
			commonx.CheckErrOrFatal(err)
			fmt.Println(res)
		} else {
			wd, err := os.Getwd()
			commonx.CheckErrOrFatal(err)
			res, err := commonx.MD5File(filepath.Join(wd, t))
			commonx.CheckErrOrFatal(err)
			fmt.Println(res)
		}
	},
}

func init() {
	rootCmd.AddCommand(md5Cmd)
}
