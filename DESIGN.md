---
name: ADP Viewer
description: Literary-journal reading surface for account development plans — true paper, deep forest, set type.
colors:
  paper: "#F6F7F4"
  surface: "#FFFFFF"
  surface-sunken: "#EEF1EC"
  ink: "#181A18"
  ink-muted: "#3A3D3A"
  ink-subtle: "#565A56"
  ink-faint: "#757A75"
  forest: "#1F3D2E"
  forest-deep: "#15301F"
  forest-tint: "#E6EFE9"
  border: "#DFE3DE"
  border-strong: "#C8CCC7"
  border-muted: "#EAEEE9"
  code-bg: "#EEF1EC"
  code-text: "#252825"
  draft-fg: "#56524E"
  draft-bg: "#ECEBE7"
  updating-fg: "#7A4A0E"
  updating-bg: "#F4E7C9"
  ready-fg: "#2F5232"
  ready-bg: "#DBE6D0"
  stale-fg: "#7A2018"
  stale-bg: "#EDD5CE"
  dark-paper: "#0E1310"
  dark-surface: "#161C18"
  dark-ink: "#E2E8E4"
  dark-forest: "#7FB994"
typography:
  display:
    fontFamily: "Fraunces, 'Iowan Old Style', Georgia, serif"
    fontWeight: 600
    letterSpacing: "-0.01em"
    lineHeight: 1.18
  body:
    fontFamily: "Newsreader, 'Iowan Old Style', Charter, Georgia, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', serif"
    fontSize: "clamp(1rem, .95rem + .2vw, 1.0625rem)"
    fontWeight: 400
    lineHeight: 1.7
  mono:
    fontFamily: "'JetBrains Mono', 'Fira Code', ui-monospace, 'SF Mono', Menlo, monospace"
    fontSize: "0.84em"
    lineHeight: 1.55
rounded:
  sm: "4px"
  md: "8px"
  lg: "12px"
spacing:
  "1": "0.25rem"
  "2": "0.5rem"
  "3": "0.75rem"
  "4": "1rem"
  "5": "1.5rem"
  "6": "2rem"
  "7": "3rem"
  "8": "4rem"
components:
  button-handoff:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    typography: "{typography.display}"
    rounded: "{rounded.sm}"
    padding: "4px 10px"
  button-handoff-hover:
    backgroundColor: "{colors.surface-sunken}"
    textColor: "{colors.ink}"
    rounded: "{rounded.sm}"
  link:
    textColor: "{colors.forest}"
    typography: "{typography.body}"
  link-hover:
    textColor: "{colors.forest-deep}"
  card:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    rounded: "{rounded.md}"
    padding: "14px"
  input-search:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    rounded: "{rounded.sm}"
    padding: "8px 12px"
  input-search-focus:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    rounded: "{rounded.sm}"
  badge-status:
    typography: "{typography.display}"
    rounded: "12px"
    padding: "2px 8px"
  tree-item-active:
    backgroundColor: "{colors.forest-tint}"
    textColor: "{colors.forest}"
  topbar:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    padding: "12px 24px"
---

# Design System: ADP Viewer

## 1. Overview

**Creative North Star: "The Reader's Desk"**

A single lamp, a set page, serious material. The ADP viewer is where an account exec sits down with the sharpest picture of an account and reads it carefully — decision chains, opportunity, the next move. The surface should feel like a well-bound book on a reader's desk, not a control panel. Penguin Classics, not PowerPoint. A literary journal, not a dashboard.

The system carries one accent — **deep forest**, the ink-like hue of a leather-bound classic — across a true off-white paper. Warmth lives in the ink and the forest green, never in a cream-tinted background. Type is the hero: Fraunces for display, Newsreader for the body, JetBrains Mono for code. Everything else recedes so the reading can come forward.

This system explicitly rejects the generic-enterprise look. If a screenshot of the viewer could be mistaken for a SaaS dashboard, a Material/Bootstrap app, an AI-tool dark-mode-with-neon, or a Notion/Confluence document shell, it has failed the brief. Distinctiveness is a principle, not a flourish.

**Key Characteristics:**
- Reading-first: long-form mixed CJK + Latin prose is the core job; every other affordance bends around legibility.
- Set, not generated: type rules, measure, leading, and hierarchy are deliberate. Nothing is default.
- One accent, used rarely: forest marks links, active state, and primary callouts only.
- Editorial voice on tool parts: file tree, badges, handoff buttons wear Fraunces and paper, not sans-serif and chrome.
- Light-first, ink-dark sibling: paper is the default; dark mode is a warm-ink companion, never a "cool tool" aesthetic.

## 2. Colors: The Forest-on-Paper Palette

One hue family carries the system. Forest (~160° in OKLCH) threads through the paper's faint tint, the primary accent, the active-state wash, and the dark theme's surfaces. Neutrals are near-true off-white, not cream.

### Primary
- **Deep Forest** (`#1F3D2E`, oklch 32% 0.06 160): the single brand accent. Used on links, active tree item text, focus rings, primary callout (note), and the theme-toggle hover border. Never used as a fill on large surfaces; its rarity is the point. On paper it hits ~11:1 contrast; on dark ink it lightens to **Fern Light** (`#7FB994`, oklch 68% 0.09 160) to stay legible against near-black.

### Secondary (semantic status)
- **Neutral Gray** (`#56524E` on `#ECEBE7`): draft — the resting, not-yet-ready state. Reads as "quiet," never as an alert.
- **Amber Ochre** (`#7A4A0E` on `#F4E7C9`): updating — in-progress, alive, warm.
- **Deep Sage** (`#2F5232` on `#DBE6D0`): ready — the only "go" signal. Kept muted, not celebration-green.
- **Oxblood Red** (`#7A2018` on `#EDD5CE`): stale — the alert. Deliberately not pure red; oxblood reads as serious, not panicky. Hue distance from the forest primary (160° vs 25°) keeps it unambiguously the alert.

### Tertiary (callout semantics)
- **Forest Note** — primary callout, shares the brand accent.
- **Sage Tip**, **Amber Warn**, **Violet Design**, **Neutral Aside** — see Components § Callouts.

### Neutral
- **Paper** (`#F6F7F4`, oklch 97% 0.004 160): the body background. True off-white with a hair of forest tint — NOT cream. The previous cream bg (`#FBF7F0`) sat in the saturated AI-default warm-neutral band; this escapes it while still reading as paper.
- **Surface** (`#FFFFFF`): cards, topbar, inputs — pure white elevated one step off the paper.
- **Surface Sunken** (`#F0EFEB`): sidebar, breadcrumb — one step below paper, anchoring chrome.
- **Ink** (`#1A1817`): headings, primary text. Warm near-black.
- **Ink Muted** (`#3A3735`): body paragraph text. Dark enough to read at ~11:1 — never lighten this "for elegance."
- **Ink Subtle** (`#56524E`): blockquotes, callout body. ~7:1.
- **Ink Faint** (`#75716C`): list markers, non-essential labels. Held at ~4.6:1 to clear WCAG AA; never lighter.
- **Border** (`#DFE3DE`), **Border Strong** (`#C8CCC7`): neutral hairlines, cool-neutral to match the forest thread.

### Named Rules
**The One Voice Rule.** Forest is the only accent. Status colors (gray/amber/sage/oxblood) are semantic signals on badges and callout borders — never decorative, never on body text, never used as a second brand color. If a screen has two accents fighting, one is wrong.

**The No-Cream Rule.** The body background is true off-white at near-zero chroma, tinted toward forest — never toward warmth by default. Warmth enters through ink, accent, and typography. A cream/sand/paper token name is a tell that the palette has drifted back to the AI default.

**The Dark-Mirror Rule.** Dark theme is the light theme mirrored through ink, not a separate "tool mode." Same hue thread (forest 160°), same semantic roles, same type. The dark accent is Fern Light (`#7FB994`), deliberately kept deeper than tech-blue range.

## 3. Typography

**Display Font:** Fraunces (with Iowan Old Style, Georgia)
**Body Font:** Newsreader (with Charter, Georgia, PingFang SC, Hiragino Sans GB, Microsoft YaHei for CJK)
**Mono Font:** JetBrains Mono (with Fira Code, SF Mono, Menlo)

**Character:** Three serifs in conversation — a high-contrast display face (Fraunces) for structure, a readable text face (Newsreader, with its italic) for prose, and a fixed-width mono for code. The pairing is unapologetically literary; it signals "sit and read" the moment the page loads. CJK characters fall through to the system Chinese serif families in Newsreader's stack, so mixed Latin + Chinese paragraphs hold a consistent texture.

### Hierarchy
- **Display / H1** (Fraunces 600, `clamp(1.7rem, 1.4rem + 1.4vw, 2.1rem)`, line-height 1.18): document title, list-page heading. Carries a hairline border beneath to set it as a section opener.
- **Headline / H2** (Fraunces 600, `clamp(1.3rem, 1.15rem + .7vw, 1.5rem)`, 1.22): major sections within a document. Also underlined with a border.
- **Title / H3** (Fraunces 600, `clamp(1.1rem, 1.03rem + .35vw, 1.2rem)`, 1.3, ink-muted): subsections. No border; recedes slightly.
- **Body** (Newsreader 400, `clamp(1rem, .95rem + .2vw, 1.0625rem)`, line-height 1.7, ink-muted): the main reading text. Measure capped at 66ch via `--measure`.
- **Label** (Fraunces 600, 11–13px, letter-spacing 0.02em): badges, buttons, topbar title, card names. The display face at small sizes — the editorial voice on tool affordances.
- **Mono** (JetBrains Mono 400, 0.84em, line-height 1.55): inline code and fenced blocks. Background `code-bg`, text `code-text`.

### Named Rules
**The Serif-Everywhere Rule.** There is no sans-serif in this system. Headings, body, labels, and buttons all carry Fraunces or Newsreader. A sans-serif appearing anywhere is a regression to the generic-tool default — refuse it.

**The Reading-Measure Rule.** Body prose is capped at 66 characters. Long lines kill the editorial voice; if a content area exceeds the measure, constrain it, don't widen the type.

**The Italic-as-Voice Rule.** Newsreader's italic is part of the palette — used in blockquotes, asides, and mermaid-raw fallback. Don't reach for bold when italic carries the emphasis more gracefully.

## 4. Elevation

Flat by default. Depth is conveyed through tonal layering (paper → surface → surface-sunken) and hairline borders, not through shadows. Shadows exist only as a quiet ambient lift on cards and a scroll-overflow cue on code blocks — never as structural elevation, never as decoration.

### Shadow Vocabulary
- **Ambient Low** (`box-shadow: 0 1px 2px rgba(18, 30, 24, .08)`): the resting lift on cards. Barely visible; its job is to separate the card from paper, not to float it.
- **Ambient Mid** (`box-shadow: 0 4px 14px rgba(18, 30, 24, .13)`): reserved for future elevated surfaces (dropdowns, dialogs). Currently unused in the resting UI.
- **Overflow Cue** (`box-shadow: inset -18px 0 14px -12px rgba(18, 30, 24, .18)`): an inset fade on the right edge of a `<pre>` that overflows its column. Signals "scroll right," nothing else.
- **Dark mirrors** (`rgba(0, 0, 0, .30)` / `.45`): the dark-theme counterparts, stronger because the surfaces are darker.

Shadow rgba values are cool-forest-neutral (`18, 30, 24`) to match the forest thread, never warm (`40, 30, 20`) — the warm shadow was a leftover from the previous ember palette.

### Named Rules
**The Flat-By-Default Rule.** Surfaces are flat at rest. Shadows appear only on cards (ambient lift) and overflow indicators (functional cue). If a shadow is decorative, remove it.

**The Tonal-Layering Rule.** Depth = tonal step (paper → surface → surface-sunken), not shadow. The sidebar is darker than paper, not floating above it.

## 5. Components

### Buttons
- **Shape:** tight 4px corners (`--radius-sm`). Editorial restraint — not pill, not square.
- **Handoff button (`.btn`):** surface background, ink text, 1px border, Fraunces 12px label, `4px 10px` padding. Hover → surface-sunken bg + forest border. Active → 1px downward nudge (`translateY(1px)`). This is the only button style; there is no primary-fill button.
- **Code copy (`.code-copy`):** surface bg, ink-subtle text, Fraunces 600 0.72rem uppercase 0.04em tracking. Hidden by default, revealed on `.code-block:hover`. Copied state → forest text. A quiet, expert affordance.
- **Theme toggle (`.theme-toggle`):** surface bg, icon + Fraunces label, 4px radius. Hover → forest border. Carries sun/moon glyphs (☀ / ☾), not text icons.

### Links
- Forest, underlined at 45% opacity, 1px thickness, 2px offset. Hover → forest-deep, full-opacity underline. Never bold, never a button-styled link.

### Cards (list page)
- **Corner:** 8px (`--radius-md`).
- **Background:** surface (`#FFFFFF`), one tonal step above paper.
- **Border:** 1px `--border`. No left-stripe accent — forbidden.
- **Shadow:** ambient-low at rest.
- **Padding:** 14px. Internal rhythm: header (name + badge) → meta row → actions row.
- **Internal hierarchy:** card name in Fraunces 600 15px as a link; meta in 12px ink-muted; handoff buttons below.

### Inputs / Search
- **Style:** 1px border, surface bg, 4px radius, 8×12 padding, 13px Newsreader.
- **Focus:** forest border + 3px forest-tinted ring (`color-mix(in srgb, forest 20%, transparent)`). The focus state is the only place forest appears as a fill-ish element.
- **No error / disabled states yet** — flag as a gap for `/impeccable harden`.

### Status Badges
- **Shape:** 12px radius (pill-ish but not full pill), 2×8 padding, Fraunces 600 11px, 0.02em tracking.
- **Four semantics:** draft (neutral gray), updating (amber ochre), ready (deep sage), stale (oxblood). See Colors § Secondary.
- **One vocabulary across list and detail:** the same badge renders identically on a card and in the topbar.

### Navigation
- **Topbar:** surface bg, 12×24 padding, bottom hairline border. Left = title or back-link; center = client name + status badge + stage; right = handoff actions + theme toggle. Fraunces throughout.
- **Breadcrumb:** surface-sunken bg, ink-muted Fraunces 13px, path segments joined by ` / `. A tonal step down from content, anchoring the reading pane.
- **File tree (sidebar):** surface-sunken bg, 280px fixed width. Items: 13px Newsreader, folder/file glyphs (📁/📄 — flagged for replacement in `/impeccable craft`), active item = forest-tint bg + forest text + 500 weight. Hover = border-muted bg.
- **No global search, no command palette** — the viewer is single-purpose. Don't add enterprise chrome.

### Callouts
- **Shape:** 8px radius, 3px left border (the only sanctioned left-stripe, because it carries semantic meaning via color), 0.9×1.1rem padding, Fraunces uppercase 0.72rem label with 0.07em tracking.
- **Five semantics:** note (forest — primary), tip (sage), warning/headsup (amber), designnote (violet), aside (transparent bg, italic body, neutral border).
- **Label always present** — a callout without a label is untyped and fails the reader.

### Markdown content
The `.content` scope carries the full prose typesetting: headings with hairline borders under H1/H2, 1.7 line-height body at 66ch measure, blockquotes (3px border-strong left, italic, ink-subtle), tables (surface-sunken header row, Fraunces headers), code blocks (code-bg, 8px radius, 1px border, overflow cue), mermaid diagrams (surface bg, centered, 8px radius).

## 6. Do's and Don'ts

### Do:
- **Do** use forest (`#1F3D2E`) as the one and only accent — on links, active tree state, focus rings, and the note callout. Its rarity is the system's voice.
- **Do** keep body paragraph text at `--ink-muted` (`#3A3D3A`) or darker. Reading text below ~7:1 on paper fails the brief.
- **Do** thread the forest hue (~160°) through every neutral — paper tint, borders, shadows, dark surfaces. Coherence is the tell that the palette was designed, not assembled.
- **Do** cap body measure at 66ch and hold line-height at 1.7 for prose. The reading experience is the product.
- **Do** use Fraunces for every label, heading, and button. The serif-on-tool-affordances is the editorial voice made structural.
- **Do** convey depth through tonal layering (paper → surface → surface-sunken), not shadows. Shadows are ambient lift only.
- **Do** run all type and color decisions through the dark mirror: if it doesn't hold at `#0E1310` paper with `#7FB994` forest, it's wrong in light too.

### Don't:
- **Don't** use a cream / sand / parchment / paper-named body background. PRODUCT.md flags the warm-neutral cream band as the saturated AI default of 2026; the previous `#FBF7F0` was squarely in it. Body bg stays true off-white at near-zero chroma, tinted toward forest.
- **Don't** let the viewer look like a generic SaaS dashboard — card grids with sidebar nav, blue primary, Inter type, the default enterprise template. PRODUCT.md names this as the #1 anti-reference.
- **Don't** use Material Design / Bootstrap patterns — elevation shadows, ripple, uniform small radius, framework-default components. No framework accent should read through.
- **Don't** ship AI-tool dark mode + neon. Pure dark + electric accent + mono everywhere reads as a toy, not a reading surface. PRODUCT.md names this directly.
- **Don't** default to the Notion / Confluence document look — sans-serif body, flat white, emoji icons, generic doc-tool aesthetic. Too flavorless for editorial content.
- **Don't** add a sans-serif anywhere. Headings, body, labels, buttons, badges all carry Fraunces or Newsreader. Sans is a regression.
- **Don't** use `border-left` greater than 1px as a decorative colored stripe on cards, list items, or alerts — the absolute-bans list. The 3px callout border is the sole sanctioned exception because it encodes semantic type.
- **Don't** apply gradient text (`background-clip: text` + gradient). Emphasis comes from weight, size, or the single forest accent — never from a gradient.
- **Don't** use glassmorphism / decorative blur. Rare and purposeful, or nothing.
- **Don't** put a tiny uppercase tracked eyebrow above every section. One deliberate kicker is voice; an eyebrow on every section is AI grammar.
- **Don't** introduce a second accent color for decoration. Status colors are semantic (gray/amber/sage/oxblood) and confined to badges + callout borders.
- **Don't** lighten `--ink-muted` for "elegance." Light muted-on-tinted-near-white is the single biggest reason designs feel hard to read.
