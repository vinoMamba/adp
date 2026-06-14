// Package skills embeds the bundled ADP skills (mirrored from adp-skill/) and
// exposes a typed catalog. The binary ships the skill *text*; it never drives
// a model — installation just writes the files into an agent's skills dir.
package skills

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/vinoMamba/adp/internal/frontmatter"
)

//go:embed data/*
var dataFS embed.FS

// Skill is one bundled ADP skill.
type Skill struct {
	Slug        string            // dir name, also the /<slug> invocation
	Name        string            // frontmatter name
	Description string            // frontmatter description
	Raw         []byte            // full SKILL.md bytes
	References  map[string][]byte // shared references/*.md, copied per skill at install
}

// All returns the bundled skills sorted by slug. Each carries the shared
// references so install can write a self-contained skill directory.
func All() ([]Skill, error) {
	refs, err := loadReferences()
	if err != nil {
		return nil, err
	}
	entries, err := fs.ReadDir(dataFS, "data")
	if err != nil {
		return nil, err
	}
	var slugs []string
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "references" {
			continue
		}
		if _, err := fs.Stat(dataFS, "data/"+e.Name()+"/SKILL.md"); err == nil {
			slugs = append(slugs, e.Name())
		}
	}
	sort.Strings(slugs)
	out := make([]Skill, 0, len(slugs))
	for _, slug := range slugs {
		raw, err := dataFS.ReadFile("data/" + slug + "/SKILL.md")
		if err != nil {
			return nil, err
		}
		name, desc := frontmatter.Parse(raw)
		out = append(out, Skill{
			Slug:        slug,
			Name:        name,
			Description: desc,
			Raw:         raw,
			References:  refs,
		})
	}
	return out, nil
}

func loadReferences() (map[string][]byte, error) {
	refs := map[string][]byte{}
	entries, err := fs.ReadDir(dataFS, "data/references")
	if err != nil {
		return nil, fmt.Errorf("read references: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		b, err := dataFS.ReadFile("data/references/" + e.Name())
		if err != nil {
			return nil, err
		}
		refs[e.Name()] = b
	}
	return refs, nil
}
