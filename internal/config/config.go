// Package config resolves filesystem locations for the ADP CLI.
package config

import (
	"os"
	"path/filepath"
)

// RootDir returns the ADP root directory holding per-client workspaces.
// Defaults to ~/Documents/adp (Windows: %USERPROFILE%\Documents\adp);
// honored override via ADP_DIR (also used by tests).
func RootDir() string {
	if d := os.Getenv("ADP_DIR"); d != "" {
		return d
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", "Documents", "adp")
	}
	return filepath.Join(home, "Documents", "adp")
}

// ConfigDir returns the global config directory (~/.adp) for config.json and
// future global state. Override via ADP_CONFIG_DIR.
func ConfigDir() string {
	if d := os.Getenv("ADP_CONFIG_DIR"); d != "" {
		return d
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".adp")
	}
	return filepath.Join(home, ".adp")
}
