package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/vinoMamba/adp/internal/store"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有客户",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := os.ReadDir(rootDir)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("根目录不存在，请先运行 adp init")
				return nil
			}
			return err
		}
		var clients []store.Client
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			if len(e.Name()) > 0 && e.Name()[0] == '.' {
				continue
			}
			c, err := store.ReadMetadata(filepath.Join(rootDir, e.Name()))
			if err != nil {
				continue
			}
			clients = append(clients, c)
		}
		sort.Slice(clients, func(i, j int) bool {
			return clients[i].Updated.After(clients[j].Updated)
		})
		if len(clients) == 0 {
			fmt.Println("暂无客户。运行 adp create <客户名称> 创建。")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tSTAGE\tSTATUS\tUPDATED")
		for _, c := range clients {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", c.Name, c.Stage, c.Status, c.Updated.Format("2006-01-02"))
		}
		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
