package serve

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gohugoio/hugo-goldmark-extensions/passthrough"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

// Chroma syntax styles, chosen to harmonize with the warm "paper"/"ember"
// palette: tango's muted browns/olives in light, gruvbox's warm ambers/oranges
// in dark. Only the syntax-token hues come from these — the code-block
// container background is owned by our --code-bg token (see pre.chroma in
// styles.css), so chroma's own background never shows through.
const (
	lightStyle = "tango"
	darkStyle  = "gruvbox"
)

// mermaidBlock matches a fenced code block whose info string is "mermaid".
var mermaidBlock = regexp.MustCompile("(?ms)^[ \t]{0,3}```[ \t]*mermaid[ \t]*\r?\n(.*?)\r?\n[ \t]{0,3}```[ \t]*$")

// calloutBlock matches a GFM-alert-style blockquote whose first line is
// `> [!TYPE]`.
var calloutBlock = regexp.MustCompile(`(?m)^[ \t]{0,3}>[ \t]*\[!(NOTE|TIP|WARNING|HEADS-UP|ASIDE|DESIGN-NOTE|PREDICT|RECALL|UNVERIFIED)\][ \t]*\r?\n((?:[ \t]{0,3}>.*(?:\r?\n|$))*)`)

// calloutLineStrip removes the `>` (and one optional following space) from the
// start of each body line of a callout, leaving the inner markdown.
var calloutLineStrip = regexp.MustCompile(`(?m)^[ \t]{0,3}> ?`)

// RenderMarkdown renders markdown to HTML with syntax highlighting, GFM
// callouts, mermaid passthrough, and LaTeX math support.
func RenderMarkdown(src []byte) ([]byte, error) {
	src = preprocessCallouts(src)
	src = preprocessMermaid(src)
	md := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle(lightStyle),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
			extension.Table,
			passthrough.New(passthrough.Config{
				InlineDelimiters: []passthrough.Delimiters{
					{Open: "$", Close: "$"},
					{Open: `\(`, Close: `\)`},
				},
				BlockDelimiters: []passthrough.Delimiters{
					{Open: "$$", Close: "$$"},
					{Open: `\[`, Close: `\]`},
				},
			}),
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// preprocessCallouts rewrites GFM-alert-style blockquotes (lines starting with
// `> [!TYPE]`) into raw <aside> HTML blocks.
func preprocessCallouts(src []byte) []byte {
	return calloutBlock.ReplaceAllFunc(src, func(match []byte) []byte {
		sub := calloutBlock.FindSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		kind := strings.ToLower(strings.ReplaceAll(string(sub[1]), "-", ""))
		label := calloutLabel(string(sub[1]))
		body := calloutLineStrip.ReplaceAll(sub[2], nil)
		var b bytes.Buffer
		b.WriteString("\n<aside class=\"callout callout-")
		b.WriteString(kind)
		b.WriteString("\">\n<p class=\"callout-label\">")
		b.WriteString(label)
		b.WriteString("</p>\n\n")
		b.Write(body)
		if !bytes.HasSuffix(body, []byte("\n")) {
			b.WriteByte('\n')
		}
		b.WriteString("\n</aside>\n\n")
		return b.Bytes()
	})
}

func calloutLabel(kind string) string {
	switch kind {
	case "DESIGN-NOTE":
		return "Design note"
	case "HEADS-UP":
		return "Heads up"
	case "NOTE":
		return "Note"
	case "TIP":
		return "Tip"
	case "WARNING":
		return "Warning"
	case "ASIDE":
		return "Aside"
	case "PREDICT":
		return "Predict"
	case "RECALL":
		return "Recall"
	case "UNVERIFIED":
		return "Unverified"
	}
	return kind
}

// preprocessMermaid rewrites ```mermaid fenced blocks into raw HTML divs that
// the browser-side mermaid library renders into SVG. The source is base64-
// encoded into a data attribute so the div has no inner content: Goldmark's
// HTML-block parser terminates a <div> at the first blank line, and mermaid
// sources routinely contain blank lines for readability — putting the text
// inline would split the diagram and double-encode the entities.
func preprocessMermaid(src []byte) []byte {
	return mermaidBlock.ReplaceAllFunc(src, func(match []byte) []byte {
		sub := mermaidBlock.FindSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		encoded := base64.StdEncoding.EncodeToString(sub[1])
		var b bytes.Buffer
		b.WriteString("\n<div class=\"mermaid\" data-encoded=\"")
		b.WriteString(encoded)
		b.WriteString("\"></div>\n")
		return b.Bytes()
	})
}

// HighlightCSS generates the chroma syntax-highlighting stylesheet, scoped per
// theme so light tokens don't leak into dark mode and vice versa.
func HighlightCSS() (template.CSS, error) {
	formatter := chromahtml.New(chromahtml.WithClasses(true))

	light := styles.Get(lightStyle)
	if light == nil {
		return "", fmt.Errorf("chroma style %q not found", lightStyle)
	}
	var lightBuf bytes.Buffer
	if err := formatter.WriteCSS(&lightBuf, light); err != nil {
		return "", err
	}

	dark := styles.Get(darkStyle)
	if dark == nil {
		return "", fmt.Errorf("chroma style %q not found", darkStyle)
	}
	var darkBuf bytes.Buffer
	if err := formatter.WriteCSS(&darkBuf, dark); err != nil {
		return "", err
	}

	var out strings.Builder
	out.WriteString(stripWrapperBackground(scopeCSS(lightBuf.String(), `:root:not([data-theme="dark"])`)))
	out.WriteString(stripWrapperBackground(scopeCSS(darkBuf.String(), `[data-theme="dark"]`)))

	return template.CSS(out.String()), nil
}

// wrapperBackground matches a single `background-color: …;` declaration, used to
// strip chroma's container background.
var wrapperBackground = regexp.MustCompile(`background-color:[^;}]*;?`)

func stripWrapperBackground(css string) string {
	var b strings.Builder
	for _, line := range strings.Split(css, "\n") {
		if strings.Contains(line, ".chroma {") || strings.Contains(line, ".bg {") {
			line = wrapperBackground.ReplaceAllString(line, "")
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

// scopeCSS prefixes every CSS rule in src with the given selector.
func scopeCSS(src, prefix string) string {
	var b strings.Builder
	for _, line := range strings.Split(src, "\n") {
		if line == "" {
			b.WriteByte('\n')
			continue
		}
		end := strings.LastIndex(line, "*/")
		if end == -1 {
			b.WriteString(line)
			b.WriteByte('\n')
			continue
		}
		b.WriteString(line[:end+2])
		b.WriteByte(' ')
		b.WriteString(prefix)
		b.WriteString(line[end+2:])
		b.WriteByte('\n')
	}
	return b.String()
}
