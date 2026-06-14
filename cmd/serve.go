package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/vinoMamba/adp/internal/serve"
)

// servePort is read by both serve and open (keep them in sync if more commands
// need the port — mirrors lathe's cmd/serve.go + cmd/open.go coupling).
var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 HTTP 服务浏览客户知识库",
	RunE: func(cmd *cobra.Command, args []string) error {
		absDir, err := filepath.Abs(rootDir)
		if err != nil {
			return fmt.Errorf("invalid directory: %w", err)
		}
		if info, err := os.Stat(absDir); err != nil || !info.IsDir() {
			return fmt.Errorf("directory does not exist: %s", absDir)
		}
		u := fmt.Sprintf("http://localhost:%d", servePort)
		openBrowser(u)
		return serve.Start(absDir, servePort)
	},
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 7260, "服务端口")
	rootCmd.AddCommand(serveCmd)
}

// openBrowser launches url in the user's default browser (best-effort).
func openBrowser(u string) {
	var c *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		c = exec.Command("open", u)
	case "linux":
		c = exec.Command("xdg-open", u)
	case "windows":
		c = exec.Command("cmd", "/c", "start", u)
	}
	if c != nil {
		c.Start() //nolint:errcheck
	}
}
