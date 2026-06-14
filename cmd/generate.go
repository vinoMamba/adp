package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateCmd is a handoff: it prints the /adp-generate skill command. The
// actual ADP regeneration happens in the user's interactive agent session.
var generateCmd = &cobra.Command{
	Use:   "generate <客户名称>",
	Short: "打印 /adp-generate 命令（在 agent 会话里粘贴以迭代 ADP 输出）",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if _, err := safeWorkspace(name); err != nil {
			return err
		}
		fmt.Printf("在你的 coding agent 会话里运行：\n\n  /adp-generate %s\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
