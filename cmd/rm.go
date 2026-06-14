package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rmForce bool

var rmCmd = &cobra.Command{
	Use:   "rm <客户名称>",
	Short: "删除客户工作区",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := safeWorkspace(args[0])
		if err != nil {
			return err
		}
		if _, err := os.Stat(workspace); os.IsNotExist(err) {
			return fmt.Errorf("客户工作区不存在：%s", workspace)
		}
		if !rmForce {
			return fmt.Errorf("将删除 %s；确认请加 --force", workspace)
		}
		if err := os.RemoveAll(workspace); err != nil {
			return fmt.Errorf("删除失败：%w", err)
		}
		fmt.Printf("已删除：%s\n", workspace)
		return nil
	},
}

func init() {
	rmCmd.Flags().BoolVar(&rmForce, "force", false, "确认删除（必需）")
	rootCmd.AddCommand(rmCmd)
}
