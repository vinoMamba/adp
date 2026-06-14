// Package cmd holds the ADP CLI subcommands (one per file).
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vinoMamba/adp/internal/config"
)

// rootDir is the ADP root directory holding per-client workspaces; resolved
// from --dir (default ~/adp).
var rootDir string

var rootCmd = &cobra.Command{
	Use:   "adp",
	Short: "ADP 客户知识库管理工具",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootDir, "dir", "d", config.RootDir(), "根目录")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
