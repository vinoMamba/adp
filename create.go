package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <客户名称>",
	Short: "创建客户目录",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		target := filepath.Join(rootDir, name)

		info, err := os.Stat(target)
		if err == nil && info.IsDir() {
			return fmt.Errorf("目录已存在: %s", target)
		}

		if err := os.MkdirAll(target, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
		fmt.Printf("已创建: %s\n", target)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
