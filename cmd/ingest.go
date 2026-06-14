package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ingestCmd is a handoff: it never drives a model. It just prints the skill
// command for the user to paste into their interactive agent session (mirrors
// lathe's lathe verify / lathe extend handoff model).
var ingestCmd = &cobra.Command{
	Use:   "ingest <客户名称>",
	Short: "打印 /adp-ingest 命令（在 agent 会话里粘贴以摄入新资料）",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if _, err := safeWorkspace(name); err != nil {
			return err
		}
		fmt.Printf("在你的 coding agent 会话里运行：\n\n  /adp-ingest %s\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ingestCmd)
}
