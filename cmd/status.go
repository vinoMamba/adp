package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vinoMamba/adp/internal/store"
)

var (
	statusStage string
	statusState string
	statusModel string
)

var statusCmd = &cobra.Command{
	Use:   "status <客户名称>",
	Short: "更新 metadata.json 的阶段 / 状态 / 模型（skill 回调）",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := safeWorkspace(args[0])
		if err != nil {
			return err
		}
		c, err := store.ReadMetadata(workspace)
		if err != nil {
			return err
		}
		if cmd.Flags().Changed("stage") {
			c.Stage = statusStage
		}
		if cmd.Flags().Changed("state") {
			s := store.Status(statusState)
			if !isValidStatus(s) {
				return fmt.Errorf("invalid --state %q (want draft|updating|ready|stale)", statusState)
			}
			c.Status = s
		}
		if cmd.Flags().Changed("model") {
			c.Model = statusModel
		}
		return store.WriteMetadata(workspace, c)
	},
}

func isValidStatus(s store.Status) bool {
	switch s {
	case store.StatusDraft, store.StatusUpdating, store.StatusReady, store.StatusStale:
		return true
	}
	return false
}

func init() {
	statusCmd.Flags().StringVar(&statusStage, "stage", "", "客户阶段")
	statusCmd.Flags().StringVar(&statusState, "state", "", "状态（draft|updating|ready|stale）")
	statusCmd.Flags().StringVar(&statusModel, "model", "", "生成 ADP 的模型标签")
	rootCmd.AddCommand(statusCmd)
}
