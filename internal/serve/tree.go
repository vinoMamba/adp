package serve

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Node is one entry in the workspace file tree (directory or .md file).
type Node struct {
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	Type     string  `json:"type"` // "file" or "dir"
	Children []*Node `json:"children,omitempty"`
}

// buildTree walks root/relPath and returns a tree of directories and .md files.
// Dotfiles and non-markdown files are skipped; empty directories are pruned.
func buildTree(root, relPath string) (*Node, error) {
	fullPath := filepath.Join(root, relPath)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	node := &Node{
		Name: filepath.Base(relPath),
		Path: relPath,
		Type: "dir",
	}

	var dirs, files []os.DirEntry
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if e.IsDir() {
			dirs = append(dirs, e)
		} else if strings.HasSuffix(strings.ToLower(name), ".md") {
			files = append(files, e)
		}
	}

	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	for _, d := range dirs {
		child, err := buildTree(root, filepath.Join(relPath, d.Name()))
		if err != nil {
			continue
		}
		if len(child.Children) == 0 {
			continue // prune empty subtrees
		}
		node.Children = append(node.Children, child)
	}

	for _, f := range files {
		node.Children = append(node.Children, &Node{
			Name: f.Name(),
			Path: filepath.Join(relPath, f.Name()),
			Type: "file",
		})
	}

	return node, nil
}
