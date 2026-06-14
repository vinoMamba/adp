package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/vinoMamba/adp/internal/skills"
)

var skillsListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出内嵌的 skills",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, err := skills.All()
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SLUG\tNAME\tDESCRIPTION")
		for _, s := range all {
			desc := s.Description
			if len(desc) > 60 {
				desc = desc[:57] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", s.Slug, s.Name, desc)
		}
		return w.Flush()
	},
}

func init() {
	skillsCmd.AddCommand(skillsListCmd)
}
