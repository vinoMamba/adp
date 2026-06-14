package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vinoMamba/adp/internal/buildinfo"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "打印版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(buildinfo.String())
	},
}

func init() {
	// Surf `adp --version` too. Set here (not in root.go) so buildinfo stays
	// colocated with the version command.
	rootCmd.Version = buildinfo.Resolve()
	rootCmd.AddCommand(versionCmd)
}
