package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed static/*
var staticFS embed.FS

var md goldmark.Markdown

func init() {
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.NewTable(),
			extension.TaskList,
			extension.Strikethrough,
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
}

func startServer(rootDir string, port int) error {
	// Serve embedded static files
	staticContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		return err
	}

	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticContent))))

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		data, err := staticFS.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	http.HandleFunc("GET /api/tree", func(w http.ResponseWriter, r *http.Request) {
		tree, err := buildTree(rootDir, "")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// Use root dir name as the top-level node name
		tree.Name = filepath.Base(rootDir)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tree)
	})

	http.HandleFunc("GET /api/file", func(w http.ResponseWriter, r *http.Request) {
		relPath := r.URL.Query().Get("p")
		if relPath == "" {
			http.Error(w, "missing p parameter", 400)
			return
		}

		// Security: prevent path traversal
		cleanPath := filepath.Clean(relPath)
		if strings.HasPrefix(cleanPath, "..") {
			http.Error(w, "invalid path", 400)
			return
		}

		fullPath := filepath.Join(rootDir, cleanPath)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("read error: %v", err), 404)
			return
		}

		var buf strings.Builder
		if err := md.Convert(data, &buf); err != nil {
			http.Error(w, fmt.Sprintf("render error: %v", err), 500)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(buf.String()))
	})

	fmt.Printf("ADP viewer serving %s on http://localhost:%d\n", rootDir, port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
