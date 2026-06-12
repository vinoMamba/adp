package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootDir string

var rootCmd = &cobra.Command{
	Use:   "adp",
	Short: "ADP 客户知识库管理工具",
}

func init() {
	defaultDir := filepath.Join(os.Getenv("HOME"), "adp")
	rootCmd.PersistentFlags().StringVarP(&rootDir, "dir", "d", defaultDir, "根目录")
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
