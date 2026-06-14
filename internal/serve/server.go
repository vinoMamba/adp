// Package serve runs the ADP viewer HTTP server: a list page of client cards
// and a per-client detail page (tree navigation + rendered markdown). Status
// badges come from each client's metadata.json; handoff buttons copy the skill
// command for the user to paste into their interactive agent session.
package serve

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/vinoMamba/adp/internal/store"
)

//go:embed list.html detail.html styles.css
var contentFS embed.FS

//go:embed static/mermaid.min.js
//go:embed static/fonts/fraunces.woff2 static/fonts/newsreader.woff2 static/fonts/newsreader-italic.woff2 static/fonts/jetbrains-mono.woff2
var staticFS embed.FS

// staticAssets whitelists the embedded files we expose under /static/.
var staticAssets = map[string]string{
	"mermaid.min.js":          "application/javascript; charset=utf-8",
	"fraunces.woff2":          "font/woff2",
	"newsreader.woff2":        "font/woff2",
	"newsreader-italic.woff2": "font/woff2",
	"jetbrains-mono.woff2":    "font/woff2",
}

// Server holds the root directory and the parsed list template.
type Server struct {
	rootDir      string
	listTmpl     *template.Template
	highlightCSS []byte
}

// cardView is the view model for one client card on the list page.
type cardView struct {
	Name           string
	Stage          string
	Status         string
	Updated        string
	MaterialsCount int
}

// Start launches the ADP viewer for rootDir on port.
func Start(rootDir string, port int) error {
	tmpl, err := template.ParseFS(contentFS, "list.html")
	if err != nil {
		return fmt.Errorf("parse list.html: %w", err)
	}
	hlCSS, err := HighlightCSS()
	if err != nil {
		return fmt.Errorf("build highlight CSS: %w", err)
	}
	s := &Server{rootDir: rootDir, listTmpl: tmpl, highlightCSS: []byte(hlCSS)}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /static/styles.css", s.handleStyles)
	mux.HandleFunc("GET /static/highlight.css", s.handleHighlightCSS)
	mux.HandleFunc("GET /static/{name}", s.handleStatic)
	mux.HandleFunc("GET /", s.handleList)
	mux.HandleFunc("GET /api/clients", s.handleAPIClients)
	mux.HandleFunc("GET /{name}", s.handleDetail)
	mux.HandleFunc("GET /{name}/api/meta", s.handleAPIMeta)
	mux.HandleFunc("GET /{name}/api/tree", s.handleAPITree)
	mux.HandleFunc("GET /{name}/api/file", s.handleAPIFile)

	fmt.Printf("ADP viewer serving %s on http://localhost:%d\n", rootDir, port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

// resolveClient maps a {name} path segment to a workspace under rootDir,
// rejecting anything that escapes the root.
func (s *Server) resolveClient(name string) (string, error) {
	clean := filepath.Clean(name)
	if clean == "" || clean == "." || clean == ".." || filepath.IsAbs(clean) {
		return "", fmt.Errorf("invalid client name")
	}
	ws := filepath.Join(s.rootDir, clean)
	rel, err := filepath.Rel(s.rootDir, ws)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("invalid client name")
	}
	return ws, nil
}

func (s *Server) gatherCards() []cardView {
	entries, err := os.ReadDir(s.rootDir)
	if err != nil {
		return nil
	}
	var cards []cardView
	for _, e := range entries {
		if !e.IsDir() || (len(e.Name()) > 0 && e.Name()[0] == '.') {
			continue
		}
		c, err := store.ReadMetadata(filepath.Join(s.rootDir, e.Name()))
		if err != nil {
			continue
		}
		cards = append(cards, cardView{
			Name:           c.Name,
			Stage:          c.Stage,
			Status:         string(c.Status),
			Updated:        c.Updated.Format("2006-01-02"),
			MaterialsCount: c.MaterialsCount,
		})
	}
	return cards
}

func (s *Server) handleStyles(w http.ResponseWriter, r *http.Request) {
	data, err := contentFS.ReadFile("styles.css")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(data)
}

func (s *Server) handleHighlightCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(s.highlightCSS)
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	contentType, ok := staticAssets[name]
	if !ok {
		http.NotFound(w, r)
		return
	}
	embedPath := "static/" + name
	if strings.HasSuffix(name, ".woff2") {
		embedPath = "static/fonts/" + name
	}
	data, err := staticFS.ReadFile(embedPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Write(data)
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	// Only the bare root serves the list; anything else is a client route.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	cards := s.gatherCards()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.listTmpl.Execute(w, map[string]any{"Clients": cards}); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (s *Server) handleAPIClients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.gatherCards())
}

func (s *Server) handleDetail(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if _, err := s.resolveClient(name); err != nil {
		http.NotFound(w, r)
		return
	}
	data, err := contentFS.ReadFile("detail.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

func (s *Server) handleAPIMeta(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	ws, err := s.resolveClient(name)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	c, err := store.ReadMetadata(ws)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (s *Server) handleAPITree(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	ws, err := s.resolveClient(name)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	tree, err := buildTree(ws, "")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	tree.Name = name
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tree)
}

func (s *Server) handleAPIFile(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	ws, err := s.resolveClient(name)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	relPath := r.URL.Query().Get("p")
	if relPath == "" {
		http.Error(w, "missing p parameter", 400)
		return
	}
	cleanPath := filepath.Clean(relPath)
	if strings.HasPrefix(cleanPath, "..") {
		http.Error(w, "invalid path", 400)
		return
	}
	fullPath := filepath.Join(ws, cleanPath)
	// Defense in depth: ensure the resolved path stays under the workspace.
	if rel, err := filepath.Rel(ws, fullPath); err != nil || strings.HasPrefix(rel, "..") {
		http.Error(w, "invalid path", 400)
		return
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("read error: %v", err), 404)
		return
	}
	out, err := RenderMarkdown(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("render error: %v", err), 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(out)
}
