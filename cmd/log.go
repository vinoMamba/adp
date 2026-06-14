package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vinoMamba/adp/internal/store"
)

var (
	logAction    string
	logJudgement string
)

var logCmd = &cobra.Command{
	Use:   "log <客户名称>",
	Short: "向 更新日志.md 追加一条记录（skill 回调）",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := safeWorkspace(args[0])
		if err != nil {
			return err
		}
		if logAction == "" {
			return cmd.Help()
		}
		if err := store.AppendUpdateLog(workspace, logAction, logJudgement); err != nil {
			return err
		}
		// CLI is the sole metadata writer: bump Updated.
		c, err := store.ReadMetadata(workspace)
		if err != nil {
			return err
		}
		return store.WriteMetadata(workspace, c)
	},
}

func init() {
	logCmd.Flags().StringVar(&logAction, "action", "", "动作（必填）")
	logCmd.Flags().StringVar(&logJudgement, "judgement", "", "判断 / 变化说明")
	rootCmd.AddCommand(logCmd)
}
