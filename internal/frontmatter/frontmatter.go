// Package frontmatter parses the name:/description: frontmatter of a SKILL.md
// without pulling in a YAML dependency.
package frontmatter

import "strings"

// Parse extracts name and description from a `---`-delimited frontmatter block
// at the top of raw. Returns ("", "") when no frontmatter is present.
func Parse(raw []byte) (name, description string) {
	lines := strings.Split(string(raw), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return "", ""
	}
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "---" {
			break
		}
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(strings.Trim(strings.TrimSpace(v), "\"'"))
		switch k {
		case "name":
			name = v
		case "description":
			description = v
		}
	}
	return name, description
}
