package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化根目录",
	RunE: func(cmd *cobra.Command, args []string) error {
		if info, err := os.Stat(rootDir); err == nil && info.IsDir() {
			fmt.Printf("目录已存在：%s\n", rootDir)
			return nil
		}
		if err := os.MkdirAll(rootDir, 0o755); err != nil {
			return fmt.Errorf("创建目录失败：%w", err)
		}
		fmt.Printf("已创建：%s\n", rootDir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
