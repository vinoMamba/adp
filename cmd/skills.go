package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vinoMamba/adp/internal/skills"
)

// skillsCmd groups the bundled-skill commands. The skills are embedded in the
// binary (internal/skills, mirrored from adp-skill/), so install works after a
// plain `go install` with no repo clone. The binary never drives a model — it
// only writes skill text into an agent's skills dir.
var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "管理内嵌的 ADP skills（安装到 Claude Code / Cursor / Codex / Gemini / opencode / Cline / Windsurf）",
}

// rawShipTarget describes an agent that consumes the raw SKILL.md verbatim.
// SKILL.md (name + description frontmatter) is now a cross-tool standard shared
// by Claude Code, Codex, Gemini CLI, opencode, Cline, and Windsurf, so these
// all ship the bytes unchanged. Cursor is the exception (needs translation) and
// keeps its own branch in installForAgent.
type rawShipTarget struct {
	display string   // human-facing name
	project []string // path segments under the project root
	user    []string // path segments under $HOME; nil/empty ⇒ project-only
}

var rawShipTargets = map[string]rawShipTarget{
	"claude-code": {display: "Claude Code", project: []string{".claude", "skills"}, user: []string{".claude", "skills"}},
	"codex":       {display: "Codex", project: []string{".agents", "skills"}, user: []string{".agents", "skills"}},
	"gemini":      {display: "Gemini", project: []string{".gemini", "skills"}, user: []string{".gemini", "skills"}},
	"opencode":    {display: "opencode", project: []string{".opencode", "skills"}, user: []string{".config", "opencode", "skills"}},
	"cline":       {display: "Cline", project: []string{".cline", "skills"}, user: []string{".cline", "skills"}},
	"windsurf":    {display: "Windsurf", project: []string{".windsurf", "skills"}},
}

// installForAgent writes every skill for one agent and returns the file count.
func installForAgent(out io.Writer, agent string, user bool, all []skills.Skill) (int, error) {
	if agent == "cursor" {
		return installCursor(out, all)
	}
	t, ok := rawShipTargets[agent]
	if !ok {
		return 0, fmt.Errorf("unknown agent %q (want one of: claude-code, cursor, codex, gemini, opencode, cline, windsurf, all)", agent)
	}
	if user && len(t.user) == 0 {
		fmt.Fprintf(out, "note: %s has no standard user-level skills dir; installing into the project instead.\n", t.display)
		user = false
	}
	dir, err := rawShipDir(t, user)
	if err != nil {
		return 0, err
	}
	fmt.Fprintf(out, "%s -> %s\n", t.display, dir)
	count := 0
	for _, s := range all {
		if err := writeRawSkill(out, dir, s); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

// writeRawSkill writes <dir>/<slug>/SKILL.md and a copy of the shared
// references/, so each installed skill is self-contained and the agent can
// resolve `references/X.md` relative to the SKILL.md.
func writeRawSkill(out io.Writer, dir string, s skills.Skill) error {
	skillDir := filepath.Join(dir, s.Slug)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", skillDir, err)
	}
	if err := writeFile(out, filepath.Join(skillDir, "SKILL.md"), s.Raw); err != nil {
		return err
	}
	for name, data := range s.References {
		refDir := filepath.Join(skillDir, "references")
		if err := os.MkdirAll(refDir, 0o755); err != nil {
			return err
		}
		if err := writeFile(out, filepath.Join(refDir, name), data); err != nil {
			return err
		}
	}
	return nil
}

// installCursor writes the Cursor slash-command translation of every skill.
// Cursor commands are single files and cannot bundle a references/ dir, so a
// note is printed — use a raw-ship agent (claude-code, codex, …) for the full
// references-backed experience.
func installCursor(out io.Writer, all []skills.Skill) (int, error) {
	fmt.Fprintln(out, "note: Cursor command files are single-file and cannot bundle references/; use claude-code/codex/gemini/opencode/cline/windsurf for the full skill.")
	dir := filepath.Join(".cursor", "commands")
	fmt.Fprintf(out, "Cursor -> %s\n", dir)
	count := 0
	for _, s := range all {
		dst := filepath.Join(dir, s.Slug+".md")
		if err := writeFile(out, dst, cursorCommand(s)); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

// cursorCommand strips the frontmatter and prepends a /<slug> header, matching
// the Cursor slash-command convention.
func cursorCommand(s skills.Skill) []byte {
	raw := string(s.Raw)
	if len(raw) >= 3 && raw[:3] == "---" {
		end := indexOf(raw, "\n---", 3)
		if end >= 0 {
			raw = raw[end+4:]
		}
	}
	return []byte("# /" + s.Slug + "\n\n" + trimLeftSpace(raw))
}

func indexOf(s, sub string, from int) int {
	for i := from; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func trimLeftSpace(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] != '\n' && s[i] != ' ' && s[i] != '\t' && s[i] != '\r' {
			return s[i:]
		}
	}
	return s
}

func rawShipDir(t rawShipTarget, user bool) (string, error) {
	if !user {
		return filepath.Join(t.project...), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(append([]string{home}, t.user...)...), nil
}

func writeFile(out io.Writer, dst string, data []byte) error {
	verb := "wrote"
	if _, err := os.Stat(dst); err == nil {
		verb = "updated"
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("mkdir for %s: %w", dst, err)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", dst, err)
	}
	fmt.Fprintf(out, "  %s %s\n", verb, dst)
	return nil
}

func init() {
	rootCmd.AddCommand(skillsCmd)
}
