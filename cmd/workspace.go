package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
)

// safeWorkspace joins rootDir with a client name and verifies the result stays
// under rootDir, rejecting empty/dotdot/absolute names. This is the path-safety
// gate shared by create/rm and any command that resolves a client by name.
func safeWorkspace(name string) (string, error) {
	clean := filepath.Clean(name)
	if clean == "" || clean == "." || clean == ".." || filepath.IsAbs(clean) {
		return "", fmt.Errorf("invalid client name: %q", name)
	}
	ws := filepath.Join(rootDir, clean)
	rel, err := filepath.Rel(rootDir, ws)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("invalid client name (escapes root): %q", name)
	}
	return ws, nil
}
