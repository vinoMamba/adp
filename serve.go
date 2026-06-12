package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 HTTP 服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		absDir, err := filepath.Abs(rootDir)
		if err != nil {
			return fmt.Errorf("invalid directory: %w", err)
		}
		if info, err := os.Stat(absDir); err != nil || !info.IsDir() {
			return fmt.Errorf("directory does not exist: %s", absDir)
		}

		url := fmt.Sprintf("http://localhost:%d", servePort)
		openBrowser(url)

		return startServer(absDir, servePort)
	},
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 7260, "服务端口")
	rootCmd.AddCommand(serveCmd)
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	}
	if cmd != nil {
		cmd.Start() //nolint:errcheck
	}
}
