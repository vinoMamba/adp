package store

import (
	"os"
	"path/filepath"
	"testing"
)

// TestScaffoldCreatesFullWorkspace locks in the contract that Scaffold produces
// every directory and template file the former Python init script did. If a
// path is added or removed, this test fails so the migration stays faithful.
func TestScaffoldCreatesFullWorkspace(t *testing.T) {
	root := t.TempDir()
	ws := filepath.Join(root, "测试公司")

	if err := Scaffold(ws, ScaffoldOptions{Name: "测试公司", Owner: "张三", Stage: "调研中"}); err != nil {
		t.Fatalf("Scaffold: %v", err)
	}

	for _, d := range workspaceDirs {
		rel := filepath.Join(ws, d)
		if info, err := os.Stat(rel); err != nil || !info.IsDir() {
			t.Errorf("missing dir %s: %v", d, err)
		}
	}

	wantFiles := []string{
		"AGENTS.md",
		"metadata.json",
		"客户知识库/索引.md",
		"客户知识库/客户画像.md",
		"客户知识库/现状.md",
		"客户知识库/人物与决策链.md",
		"客户知识库/机会与动机.md",
		"客户知识库/行动计划.md",
		"客户知识库/来源登记.md",
		"客户知识库/更新日志.md",
		"输出/测试公司-ADP.md",
	}
	for _, f := range wantFiles {
		if _, err := os.Stat(filepath.Join(ws, f)); err != nil {
			t.Errorf("missing file %s: %v", f, err)
		}
	}

	// ADP output has the client name substituted.
	adp, err := os.ReadFile(filepath.Join(ws, "输出/测试公司-ADP.md"))
	if err != nil {
		t.Fatalf("read ADP output: %v", err)
	}
	if got := string(adp); got[:len("# 测试公司")] != "# 测试公司" {
		t.Errorf("ADP output head = %q, want # 测试公司", got[:12])
	}

	// metadata.json round-trips with the scaffolded values.
	c, err := ReadMetadata(ws)
	if err != nil {
		t.Fatalf("ReadMetadata: %v", err)
	}
	if c.Name != "测试公司" || c.Owner != "张三" || c.Stage != "调研中" || c.Status != StatusDraft {
		t.Errorf("metadata = %+v", c)
	}
}

// TestScaffoldIdempotentGuard verifies re-scaffolding an existing workspace is
// refused (the 客户知识库/索引.md sentinel), matching the Python default.
func TestScaffoldIdempotentGuard(t *testing.T) {
	ws := filepath.Join(t.TempDir(), "客户A")
	if err := Scaffold(ws, ScaffoldOptions{Name: "客户A"}); err != nil {
		t.Fatalf("first Scaffold: %v", err)
	}
	if err := Scaffold(ws, ScaffoldOptions{Name: "客户A"}); err == nil {
		t.Fatal("second Scaffold should fail on existing workspace")
	}
}

// TestAppendUpdateLog covers both the new-section and append-to-today paths.
func TestAppendUpdateLog(t *testing.T) {
	ws := filepath.Join(t.TempDir(), "客户B")
	if err := Scaffold(ws, ScaffoldOptions{Name: "客户B"}); err != nil {
		t.Fatalf("Scaffold: %v", err)
	}
	if err := AppendUpdateLog(ws, "摄入", "补充纪要"); err != nil {
		t.Fatalf("AppendUpdateLog: %v", err)
	}
	if err := AppendUpdateLog(ws, "生成", "首版"); err != nil {
		t.Fatalf("AppendUpdateLog 2: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(ws, "客户知识库/更新日志.md"))
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if c := string(got); !contains(c, "摄入：补充纪要") || !contains(c, "生成：首版") {
		t.Errorf("log missing entries:\n%s", got)
	}
}

func contains(s, sub string) bool { return len(s) >= len(sub) && (s == sub || indexOfStr(s, sub) >= 0) }
func indexOfStr(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
