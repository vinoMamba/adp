# Product

## Register

product

> Default is `product` for the current viewer surface. A future brand surface (marketing landing, public docs for the CLI itself) may switch to `brand` — treat that as a separate task, not a re-skin of the viewer.

## Users

**Primary:** Enterprise account executives / BD preparing for or reviewing client meetings. Chinese-speaking, operating in B2B enterprise sales.

**Context of use:** Opens the viewer on a laptop, often mid-call or between meetings. Two reading modes, both must be excellent:

- **Scan mode** — answer "what's the status of this account, what's the next move, is anything stale?" in under 5 seconds. Status badges, stage, last-updated, materials count must read at a glance.
- **Read mode** — drop into the ADP document itself and read it carefully: decision chains, opportunity analysis, action plans. Long-form Chinese + English mixed prose, tables, callouts, diagrams.

**Job to be done:** Walk into (or out of) a client conversation with the sharpest possible picture of the account — who matters, what's happening, what to do next — without re-reading everything.

## Product Purpose

ADP (Account Development Plan / 客户经营计划) is a CLI + coding-agent skills toolkit that maintains living account development plans for enterprise accounts.

- **The Go CLI** owns durable state: workspace scaffolding, `metadata.json`, `更新日志.md`, `来源登记.md`, and the viewer HTTP server. It never drives a model.
- **The skills** (`/adp`, `/adp-ingest`, `/adp-generate`, `/adp-ask`, `/adp-review`) run in the interactive coding-agent session and do all LLM work, calling back into the CLI to record state.
- **The viewer** (`adp serve`) is the design surface this document governs: a list page (client cards) and a detail page (file-tree sidebar + rendered markdown).

Success = an account exec treats the viewer as the single source of truth for the account and trusts it enough to bring into a meeting. If they re-open their old notes doc instead, it failed.

## Brand Personality

**Editorial book, not enterprise dashboard.** Three words: *considered, warm, expert.*

- **Considered** — every type rule, spacing choice, and color has a reason. Nothing is default. The page feels set, not generated.
- **Warm** — paper, ember, ink. A reading surface you'd sit with, not a control panel you'd operate. Warmth comes from palette, type, and measure — never from cream-bg-by-default.
- **Expert** — the voice of someone who has done this account work for years and knows what matters. Confident hierarchy, no hedging, no decoration that doesn't earn its place.

Voice carries into copy: status labels and handoff buttons are plain and imperative (摄入 / 生成 / 提问 / 审计), not cute.

## Anti-references

Every default enterprise look is on this list. If the viewer could be mistaken for any of these, it has failed the brief.

- **Generic SaaS dashboard** — card grids with sidebar nav, blue primary, Inter type, the default enterprise template. The #1 thing to not look like.
- **Material Design / Bootstrap** — elevation shadows, ripple, uniform small radius, the framework-default look. No framework accent should read through.
- **AI-tool dark mode + neon** — pure dark + electric accent + mono everywhere. Reads as a toy, not a reading surface.
- **Notion / Confluence document app** — sans-serif body, flat white, emoji icons, the generic "doc tool" aesthetic. Too flavorless for editorial content.

The tell: if a designer looking at a screenshot couldn't guess "this is the ADP viewer," it's too generic.

## Design Principles

1. **Reading experience first.** The viewer renders long-form mixed CJK + Latin markdown that people read carefully. Optimize for sustained reading: measure, leading, hierarchy, contrast, paragraph rhythm. Every other concern (status scanning, handoff affordances) bends around this, never the reverse.

2. **Editorial voice, tool affordances.** The editorial-book aesthetic is the brand. Tool parts — file tree, status badges, copy-to-handoff buttons — wear the editorial voice (serif labels, paper surfaces, considered motion), not the other way around. Never let a SaaS component pattern override the typographic system.

3. **Distinctive by default.** Every anti-reference above is a default. The path of least resistance (blue primary, sans body, card grid, framework components) is the failure mode. When a decision is close, pick the option that makes the result less mistakable for a generic tool.

4. **Scan, then read.** Both modes must be excellent. Status, stage, staleness, and next-action signaling must resolve in under 5 seconds of list-page scanning. Then the drop-into-the-document experience must reward the deep read. Neither mode is optional; neither mode is the "real" one.

5. **Practice what it preaches.** The ADP itself is a client-facing document advocating considered account work. The viewer's craft should model that craft — a viewer that looks generic undermines the methodology it displays.

## Accessibility & Inclusion

**Target: WCAG 2.1 AA.**

- **Contrast** — body text ≥ 4.5:1 against bg; large/bold text ≥ 3:1; placeholder text meets 4.5:1, not the muted-gray default. The warm paper palette makes this non-obvious — verify every muted tone.
- **Motion** — respect `prefers-reduced-motion`. Every animation has a still or crossfade alternative.
- **Keyboard** — full keyboard navigation: search, tree, links, copy buttons, theme toggle. Visible focus rings.
- **Screen readers** — semantic HTML, aria-labels on icon-only controls (theme toggle, copy), status badges with accessible names.
- **CJK + Latin** — mixed-script rendering must stay legible at every breakpoint. Test Chinese line-breaking and Latin/Chinese spacing; don't assume Latin-only metrics.

## Tech stack (for future commands)

- Go HTTP server (`internal/serve/`), no JS framework.
- Vanilla HTML templates (`list.html`, `detail.html`) + single `styles.css`.
- Self-hosted fonts: Fraunces (display), Newsreader (body + italic), JetBrains Mono.
- goldmark for markdown + chroma highlighting + mermaid.js for diagrams.
- Light/dark theme via `[data-theme]` on `<html>`, persisted in localStorage.
