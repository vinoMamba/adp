package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/vinoMamba/adp/internal/store"
)

var (
	srcType      string
	srcDate      string
	srcAuthority string
	srcPage      string
	srcKey       string
	srcNote      string
)

var sourceCmd = &cobra.Command{
	Use:   "source <客户名称> --origin <来源>",
	Short: "向 来源登记.md 追加一行（skill 回调）",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := safeWorkspace(args[0])
		if err != nil {
			return err
		}
		origin, _ := cmd.Flags().GetString("origin")
		if origin == "" {
			return cmd.Help()
		}
		date := srcDate
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}
		entry := store.SourceEntry{
			Origin:    origin,
			Type:      srcType,
			Date:      date,
			Authority: srcAuthority,
			Page:      srcPage,
			KeyFields: srcKey,
			Note:      srcNote,
		}
		if err := store.AppendSource(workspace, entry); err != nil {
			return err
		}
		c, err := store.ReadMetadata(workspace)
		if err != nil {
			return err
		}
		c.MaterialsCount++
		return store.WriteMetadata(workspace, c)
	},
}

func init() {
	sourceCmd.Flags().String("origin", "", "来源（必填）")
	sourceCmd.Flags().StringVar(&srcType, "type", "", "类型")
	sourceCmd.Flags().StringVar(&srcDate, "date", "", "日期（默认今天）")
	sourceCmd.Flags().StringVar(&srcAuthority, "authority", "", "权威层级")
	sourceCmd.Flags().StringVar(&srcPage, "page", "", "影响页面")
	sourceCmd.Flags().StringVar(&srcKey, "key", "", "关键字段 / 数字")
	sourceCmd.Flags().StringVar(&srcNote, "note", "", "备注")
	rootCmd.AddCommand(sourceCmd)
}
