package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vinoMamba/adp/internal/skills"
)

var (
	installAgent string
	installUser  bool
)

var skillsInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "把内嵌的 skills 写进 agent 的 skills 目录",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, err := skills.All()
		if err != nil {
			return err
		}
		if len(all) == 0 {
			return fmt.Errorf("no skills bundled")
		}
		agents := []string{installAgent}
		if installAgent == "all" {
			agents = []string{"claude-code", "cursor", "codex", "gemini", "opencode", "cline", "windsurf"}
		}
		total := 0
		for _, a := range agents {
			n, err := installForAgent(cmd.OutOrStdout(), a, installUser, all)
			if err != nil {
				return err
			}
			total += n
		}
		fmt.Printf("\n%d skill(s) installed.\n", len(all))
		return nil
	},
}

func init() {
	skillsInstallCmd.Flags().StringVar(&installAgent, "agent", "claude-code", "目标 agent：claude-code|cursor|codex|gemini|opencode|cline|windsurf|all")
	skillsInstallCmd.Flags().BoolVar(&installUser, "user", false, "安装到用户级目录（~）而非项目级（.）")
	skillsCmd.AddCommand(skillsInstallCmd)
}
