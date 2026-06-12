package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Node struct {
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	Type     string  `json:"type"` // "file" or "dir"
	Children []*Node `json:"children,omitempty"`
}

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

	// Collect valid entries first for stable sorting
	var dirs, files []os.DirEntry
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") || name == ".DS_Store" {
			continue
		}
		if e.IsDir() {
			dirs = append(dirs, e)
		} else if strings.HasSuffix(strings.ToLower(name), ".md") {
			files = append(files, e)
		}
	}

	// Sort dirs and files alphabetically
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	// Add dirs first, then files
	for _, d := range dirs {
		childPath := filepath.Join(relPath, d.Name())
		child, err := buildTree(root, childPath)
		if err != nil {
			continue
		}
		// Skip empty directories (no md files in subtree)
		if len(child.Children) == 0 {
			continue
		}
		node.Children = append(node.Children, child)
	}

	for _, f := range files {
		childPath := filepath.Join(relPath, f.Name())
		node.Children = append(node.Children, &Node{
			Name: f.Name(),
			Path: childPath,
			Type: "file",
		})
	}

	return node, nil
}
