package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ScaffoldOptions configures Scaffold.
type ScaffoldOptions struct {
	Name  string // client name (also the workspace leaf dir)
	Owner string
	Stage string
}

// workspaceDirs are the directories created in every client workspace.
var workspaceDirs = []string{
	"原始资料/公开调研",
	"原始资料/拜访纪要",
	"原始资料/方案报价",
	"原始资料/CRM记录",
	"原始资料/系统资料",
	"客户知识库",
	"输出",
}

// Scaffold creates a full client ADP workspace at workspace: the directory
// tree, knowledge-base template files, the initial ADP output, and a
// metadata.json. This is the Go port of the former
// adp-skill/scripts/init_adp_workspace.py.
//
// Idempotency: if 客户知识库/索引.md already exists, Scaffold returns an
// error (matching the Python default of refusing to overwrite an initialized
// workspace). Individual files are written only when absent, so a partially
// populated workspace is completed rather than clobbered.
func Scaffold(workspace string, opt ScaffoldOptions) error {
	indexPath := filepath.Join(workspace, "客户知识库", "索引.md")
	if _, err := os.Stat(indexPath); err == nil {
		return fmt.Errorf("客户 ADP 工作区已存在：%s", workspace)
	}

	for _, d := range workspaceDirs {
		if err := os.MkdirAll(filepath.Join(workspace, d), 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", d, err)
		}
	}

	today := time.Now().Format("2006-01-02")
	for rel, content := range templateFiles(opt.Name, opt.Owner, opt.Stage, today) {
		full := filepath.Join(workspace, rel)
		if _, err := os.Stat(full); err == nil {
			continue // preserve existing file
		}
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return fmt.Errorf("mkdir for %s: %w", rel, err)
		}
		body := strings.TrimRight(strings.TrimSpace(content), "\n") + "\n"
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", rel, err)
		}
	}

	meta := Client{
		Name:    opt.Name,
		Owner:   opt.Owner,
		Stage:   opt.Stage,
		Status:  StatusDraft,
		Created: time.Now().UTC(),
	}
	return WriteMetadata(workspace, meta)
}
