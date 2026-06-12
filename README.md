# ADP

A client knowledge base CLI tool. Browse local Markdown documents with a built-in HTTP server.

## Install

```bash
make install
```

## Usage

```bash
# Initialize root directory (default ~/adp)
adp init

# Create client directory
adp create ClientA
adp create TestCompany

# Start server, opens browser automatically
adp serve

# Custom port or directory
adp serve -p 9090 -d /path/to/docs
```

## Commands

| Command | Description |
|---------|-------------|
| `adp init` | Initialize root directory |
| `adp create <name>` | Create client subdirectory |
| `adp serve` | Start HTTP server to browse documents |

## Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-d, --dir` | `~/adp` | Root directory path |

## Serve Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-p, --port` | `7260` | HTTP server port |

## Directory Structure

```
~/adp/
├── ClientA/
├── ClientB/
└── ...
```

Place `.md` files under each client directory. The serve command renders them as HTML with a tree navigation sidebar.
