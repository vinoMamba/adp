package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AppendUpdateLog appends a dated entry to 客户知识库/更新日志.md. If today's
// `## <date>` section already exists the entry is appended to it; otherwise a
// new section is inserted at the top (newest first). The file is created from
// the canonical header if missing. This does NOT touch metadata.json — the CLI
// layer (cmd) owns metadata writes.
func AppendUpdateLog(workspace, action, judgement string) error {
	path := filepath.Join(workspace, "客户知识库", "更新日志.md")
	content, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		content = []byte("# 更新日志\n")
	}
	today := time.Now().Format("2006-01-02")
	entry := strings.TrimSpace(action)
	if j := strings.TrimSpace(judgement); j != "" {
		entry += "：" + j
	}
	updated := insertDatedEntry(string(content), today, entry)
	return os.WriteFile(path, []byte(updated), 0o644)
}

// insertDatedEntry places `- entry` under the `## today` section of a markdown
// log, creating the section at the top (right after the `# ` title) when absent.
func insertDatedEntry(content, today, entry string) string {
	content = strings.TrimRight(content, "\n") + "\n"
	lines := strings.Split(content, "\n")

	// Locate the `# title` header end.
	headerEnd := 0
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "# ") {
			headerEnd = i + 1
			break
		}
	}

	todayHeader := "## " + today
	// Locate today's section.
	for i := headerEnd; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == todayHeader {
			// Append before the next `## ` or EOF.
			sectionEnd := len(lines)
			for j := i + 1; j < len(lines); j++ {
				if strings.HasPrefix(strings.TrimSpace(lines[j]), "## ") {
					sectionEnd = j
					break
				}
			}
			insertAt := sectionEnd
			for j := sectionEnd - 1; j > i; j-- {
				if strings.TrimSpace(lines[j]) != "" {
					insertAt = j + 1
					break
				}
			}
			out := append([]string{}, lines[:insertAt]...)
			out = append(out, "- "+entry)
			out = append(out, lines[insertAt:]...)
			return strings.TrimRight(strings.Join(out, "\n"), "\n") + "\n"
		}
	}

	// New section at the top.
	section := []string{"", todayHeader, "", "- " + entry, ""}
	out := append([]string{}, lines[:headerEnd]...)
	out = append(out, section...)
	out = append(out, lines[headerEnd:]...)
	return strings.TrimRight(strings.Join(out, "\n"), "\n") + "\n"
}

// SourceEntry is one row of 客户知识库/来源登记.md.
type SourceEntry struct {
	Origin    string // 来源
	Type      string // 类型
	Date      string // 日期
	Authority string // 权威层级
	Page      string // 影响页面
	KeyFields string // 关键字段 / 数字
	Note      string // 备注
}

// AppendSource appends a row to 客户知识库/来源登记.md, creating the table
// header if the file is missing.
func AppendSource(workspace string, e SourceEntry) error {
	path := filepath.Join(workspace, "客户知识库", "来源登记.md")
	content, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		content = []byte("# 来源登记\n\n| 来源 | 类型 | 日期 | 权威层级 | 影响页面 | 关键字段 / 数字 | 备注 |\n|---|---|---|---|---|---|---|\n")
	}
	row := fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |",
		e.Origin, e.Type, e.Date, e.Authority, e.Page, e.KeyFields, e.Note)
	updated := strings.TrimRight(string(content), "\n") + "\n" + row + "\n"
	return os.WriteFile(path, []byte(updated), 0o644)
}
