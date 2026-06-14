# ADP

A client knowledge base CLI + a set of coding-agent skills that together maintain Account Development Plans (客户经营计划) for enterprise accounts.

The boundary is the same one lathe proved out: **skills generate content; the CLI owns durable state.** The Go binary never drives a model — all LLM work runs in your interactive coding-agent session via user-invoked skills. The CLI scaffolds workspaces, owns `metadata.json` / `更新日志.md` / `来源登记.md`, and serves the viewer UI; skills call back into the CLI to record state.

## Install

### One-liner (recommended)

Picks the latest release for your platform, verifies SHA256 against the published `checksums.txt`, and installs `adp` to `~/.local/bin`:

```bash
curl -fsSL https://github.com/vinoMamba/adp/releases/latest/download/adp-installer.sh | bash
```

Pin a version with `... | bash -s -- --version v0.1.1`, or install elsewhere with `... | bash -s -- --prefix /usr/local/bin`.

### Build from source

```bash
make install          # builds and copies `adp` to ~/.local/bin
```

### Update later

```bash
adp update            # self-update from GitHub Releases (SHA256-verified, atomic swap)
```

### For AI agents

If you are an AI coding agent installing `adp` on the user's machine, run the explicit steps below instead of piping to bash — each step is independently verifiable, and there is no `curl | bash` trust hop. Copy-paste as one block; the last line is the success criterion.

```bash
set -euo pipefail
# 1. detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$OS/$ARCH" in
  darwin/arm64|darwin/aarch64) asset=adp_darwin_arm64.tar.gz ;;
  darwin/x86_64|darwin/amd64)  asset=adp_darwin_amd64.tar.gz ;;
  linux/arm64|linux/aarch64)   asset=adp_linux_arm64.tar.gz ;;
  linux/x86_64|linux/amd64)    asset=adp_linux_amd64.tar.gz ;;
  *) echo "unsupported: $OS/$ARCH" >&2; exit 1 ;;
esac
# 2. resolve latest release
tag=$(curl -fsSL https://api.github.com/repos/vinoMamba/adp/releases/latest \
      | sed -nE 's/.*"tag_name"[[:space:]]*:[[:space:]]*"(v[^"]+)".*/\1/p' | head -1)
ver=${tag#v}
base="https://github.com/vinoMamba/adp/releases/download/$tag"
# 3. download archive + checksums to a temp dir
tmp=$(mktemp -d); trap 'rm -rf "$tmp"' EXIT
curl -fsSL -o "$tmp/$asset"        "$base/$asset"
curl -fsSL -o "$tmp/checksums.txt" "$base/adp_${ver}_checksums.txt"
# 4. verify SHA256
want=$(awk -v f="$asset" '$2==f {print $1}' "$tmp/checksums.txt")
have=$(sha256sum "$tmp/$asset" 2>/dev/null | awk '{print $1}' \
       || shasum -a 256 "$tmp/$asset" | awk '{print $1}')
[ "$(printf %s "$want" | tr A-Z a-z)" = "$(printf %s "$have" | tr A-Z a-z)" ] \
  || { echo "checksum mismatch for $asset" >&2; exit 1; }
# 5. extract + install
tar -xzf "$tmp/$asset" -C "$tmp" adp
mkdir -p "$HOME/.local/bin"
install -m 0755 "$tmp/adp" "$HOME/.local/bin/adp"
# 6. verify (must print the resolved tag)
export PATH="$HOME/.local/bin:$PATH"
adp version
```

Expected: the last command prints `adp <tag>` matching the resolved `$tag`. If `~/.local/bin` is not on the user's PATH, suggest adding `export PATH="$HOME/.local/bin:$PATH"` to their shell profile. Windows is not supported by this snippet — point the user at the `.zip` on the [releases page](https://github.com/vinoMamba/adp/releases/latest).

## The skills (run in your agent session)

After installing, drop them into a project so your agent discovers them:

```bash
adp skills install                  # ./.claude/skills/<name>/SKILL.md
adp skills install --user           # ~/.claude/skills/<name>/SKILL.md
adp skills install --agent cursor   # ./.cursor/commands/<slug>.md
adp skills install --agent codex    # ./.agents/skills/<name>/SKILL.md
adp skills install --agent gemini   # ./.gemini/skills/<name>/SKILL.md
adp skills install --agent all      # every target above
adp skills list
```

Five skills, one per intent. Each is self-contained (ships its own copy of `references/`):

| Skill | Invocation | Does |
|---|---|---|
| `adp` | `/adp <客户名称>` | Scaffold workspace (`adp create`), ingest any provided materials, generate the first ADP |
| `adp-ingest` | `/adp-ingest <客户名称>` | Ingest new raw materials into the knowledge pages |
| `adp-generate` | `/adp-generate <客户名称>` | Iterate the standard 10-section ADP output from the knowledge base |
| `adp-ask` | `/adp-ask <客户名称>` | Answer questions about a client (read-only) |
| `adp-review` | `/adp-review <客户名称>` | Audit the ADP against the quality gates |

## CLI commands

| Command | Description |
|---|---|
| `adp init` | Initialize root directory (default `~/Documents/adp`) |
| `adp create <客户名称> [--owner X] [--stage Y]` | Create a full client workspace (dirs + templates + metadata.json) |
| `adp list` | List all clients with stage / status / updated |
| `adp open [客户名称]` | Open the viewer in the browser |
| `adp rm <客户名称> --force` | Remove a client workspace |
| `adp serve [-p 7260]` | Start the viewer HTTP server |
| `adp ingest <客户名称>` | Print the `/adp-ingest` command to paste (handoff) |
| `adp generate <客户名称>` | Print the `/adp-generate` command to paste (handoff) |
| `adp log <客户名称> --action --judgement` | Append to 更新日志.md (skill callback) |
| `adp source <客户名称> --origin --type [...]` | Append to 来源登记.md (skill callback) |
| `adp status <客户名称> [--stage] [--state] [--model]` | Update metadata.json (skill callback) |
| `adp skills install/list` | Manage bundled skills |
| `adp update [--check] [--version v0.2.0]` | Self-update from GitHub Releases (SHA256-verified, atomic swap) |
| `adp version` | Print version info |

## Global flags

| Flag | Default | Description |
|---|---|---|
| `-d, --dir` | `~/Documents/adp` | Root directory (override via `ADP_DIR`) |

## Directory structure

```
~/Documents/adp/
├── <客户名称>/
│   ├── metadata.json              # English-keyed client metadata
│   ├── AGENTS.md
│   ├── 原始资料/{公开调研,拜访纪要,方案报价,CRM记录,系统资料}/
│   ├── 客户知识库/{索引,客户画像,现状,人物与决策链,机会与动机,行动计划,来源登记,更新日志}.md
│   └── 输出/<客户名称>-ADP.md      # the single, iterated ADP output
└── ...
```

## metadata.json

```json
{
  "name": "客户名称",
  "owner": "负责人",
  "stage": "调研中",
  "status": "draft",
  "created": "...",
  "updated": "...",
  "model": "Claude Opus 4.8",
  "materials_count": 0
}
```

`status`: `draft` → `updating` → `ready`; `stale` when new materials arrive but the ADP isn't regenerated. Field keys are English; values may be Chinese.

## Development

```bash
make build          # go build
make check          # pre-PR gate: skillsCheck + vet + test
make skills         # sync adp-skill/ → internal/skills/data/ (the embed mirror)
make skillsCheck    # fail if the mirror drifted
```

`adp-skill/` is the single source of truth for skills. Because `go:embed` cannot reach it from `internal/skills/`, a tracked mirror at `internal/skills/data/` is the embed source; `make skills` regenerates it. Never hand-edit the mirror — edit `adp-skill/` and run `make skills`.
