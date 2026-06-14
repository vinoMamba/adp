package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [客户名称]",
	Short: "在浏览器打开知识库",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		u := fmt.Sprintf("http://localhost:%d/", servePort)
		if len(args) == 1 {
			u += url.PathEscape(args[0])
		}
		openBrowser(u)
		fmt.Println(u)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
