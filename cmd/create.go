package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vinoMamba/adp/internal/store"
)

var (
	createOwner string
	createStage string
)

var createCmd = &cobra.Command{
	Use:   "create <客户名称>",
	Short: "创建客户工作区（完整脚手架 + metadata）",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		workspace, err := safeWorkspace(name)
		if err != nil {
			return err
		}
		opt := store.ScaffoldOptions{
			Name:  name,
			Owner: createOwner,
			Stage: createStage,
		}
		if err := store.Scaffold(workspace, opt); err != nil {
			return err
		}
		fmt.Printf("已初始化客户 ADP 工作区：%s\n", workspace)
		return nil
	},
}

func init() {
	createCmd.Flags().StringVar(&createOwner, "owner", "待填写", "客户负责人")
	createCmd.Flags().StringVar(&createStage, "stage", "待确认", "客户当前阶段")
	rootCmd.AddCommand(createCmd)
}
