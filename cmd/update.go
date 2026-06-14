package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/vinoMamba/adp/internal/buildinfo"
	"github.com/vinoMamba/adp/internal/updater"
)

var (
	updateCheck   bool
	updateVersion string
	updateForce   bool
)

// updateCmd implements `adp update`: query the latest GitHub release and
// atomically replace the running executable when a newer version exists.
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "检查并更新到最新版本（从 GitHub Releases）",
	Long: `查询 GitHub Releases 上的最新版本，与当前版本比较，
若有更新则下载对应平台的二进制（含 SHA256 校验）并原子替换当前可执行文件。

  adp update              # 检查并更新到最新版本
  adp update --check      # 只检查不更新
  adp update --version v0.2.0   # 更新到指定版本
  adp update --force      # 即使已是最新也重新下载替换`,
	RunE: runUpdate,
}

func init() {
	updateCmd.Flags().BoolVarP(&updateCheck, "check", "c", false, "只检查是否有新版本，不执行更新")
	updateCmd.Flags().StringVarP(&updateVersion, "version", "v", "", "更新到指定版本（如 v0.2.0），默认为最新")
	updateCmd.Flags().BoolVar(&updateForce, "force", false, "即使已是最新版本也强制重新下载替换")
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, _ []string) error {
	out := cmd.OutOrStdout()
	current := buildinfo.Resolve()
	currentLabel := current
	if currentLabel == "" {
		currentLabel = "dev (未发布版本)"
	}

	// Best-effort: clean up any stale .old binary on Windows from a prior run.
	updater.CleanupStaleWindowsBackup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	fmt.Fprintf(out, "当前版本: %s\n", currentLabel)
	fmt.Fprintf(out, "正在查询最新版本…\n")

	client := &http.Client{Timeout: 60 * time.Second}
	rel, err := updater.Latest(ctx, client, buildinfo.Repo, updateVersion)
	if err != nil {
		if errors.Is(err, updater.ErrNotFound) {
			if updateVersion == "" {
				return fmt.Errorf("没有可用的发布版本（仓库 %s 尚未发布 release）", buildinfo.Repo)
			}
			return fmt.Errorf("找不到版本 %s：请确认 tag 存在", updateVersion)
		}
		return fmt.Errorf("查询版本失败：%w", err)
	}

	fmt.Fprintf(out, "最新版本: %s", rel.Tag)
	if !rel.Published.IsZero() {
		fmt.Fprintf(out, "  (发布于 %s)", rel.Published.Format("2006-01-02"))
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, rel.HTMLURL)

	if updateCheck {
		// --check 只报告状态，不做改动。
		if updater.NeedUpdate(current, rel.Tag) {
			fmt.Fprintf(out, "\n有新版本可用：%s → %s\n", currentLabel, rel.Tag)
			fmt.Fprintln(out, "重新运行 `adp update` 来更新。")
			return nil
		}
		fmt.Fprintln(out, "\n已是最新版本。")
		return nil
	}

	if current == "" && !updateForce {
		// dev 构建：缺少版本基线，无法判断"是否需要更新"。
		// 提示用户用 --force 显式确认，避免误把本地构建覆盖成 release 版本。
		return fmt.Errorf("当前为 dev 构建，无法判断版本。\n如需强制安装 %s，请用 `adp update --force`", rel.Tag)
	}

	if !updater.NeedUpdate(current, rel.Tag) && !updateForce {
		fmt.Fprintln(out, "\n已是最新版本。")
		return nil
	}

	arrow := "→"
	if current == "" {
		arrow = "（dev）→"
	}
	fmt.Fprintf(out, "\n更新中：%s %s %s\n", currentLabel, arrow, rel.Tag)
	fmt.Fprintf(out, "平台: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	fmt.Fprint(out, "下载并校验中…\n")
	body, err := rel.Download(ctx, client, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return fmt.Errorf("下载失败：%w", err)
	}
	fmt.Fprintf(out, "已下载 %s 并通过校验\n", humanBytes(int64(len(body))))

	// Apply 替换 os.Executable() 指向的二进制。
	target, _ := os.Executable()
	fmt.Fprintf(out, "替换 %s…\n", target)
	if err := updater.Apply("", body); err != nil {
		return fmt.Errorf("替换二进制失败：%w", err)
	}

	fmt.Fprintf(out, "\n更新完成：%s %s %s\n", currentLabel, arrow, rel.Tag)
	if rel.HTMLURL != "" {
		fmt.Fprintln(out, "查看更新内容：", rel.HTMLURL)
	}
	return nil
}

// humanBytes formats a byte count with a sensible unit. Used for the
// "downloaded N" status line; not for general-purpose display.
func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for x := n / unit; x >= unit; x /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
}
