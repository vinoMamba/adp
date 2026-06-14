#!/usr/bin/env bash
#
# adp installer: downloads the latest release binary for the current
# platform, verifies SHA256 against the goreleaser-published checksums.txt,
# extracts, and installs `adp` to ${INSTALL_DIR:-$HOME/.local/bin}.
#
# Usage:
#   curl -fsSL https://github.com/vinoMamba/adp/releases/latest/download/adp-installer.sh | bash
#   INSTALL_DIR=/usr/local/bin bash adp-installer.sh
#   adp-installer.sh --version v0.1.1   # pin a specific release
#
# Exit codes: 0 ok · 1 generic · 2 unsupported platform · 3 checksum mismatch
set -euo pipefail

REPO="vinoMamba/adp"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
PIN_VERSION=""

while [ $# -gt 0 ]; do
  case "$1" in
    --version) PIN_VERSION="$2"; shift 2 ;;
    --prefix)  INSTALL_DIR="$2"; shift 2 ;;
    -h|--help)
      sed -n '2,11p' "$0" | sed 's/^# \{0,1\}//'
      exit 0 ;;
    *) echo "unknown arg: $1" >&2; exit 1 ;;
  esac
done

# --- 1. platform detect ---
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$OS/$ARCH" in
  darwin/arm64|darwin/aarch64) OS=darwin; ARCH=arm64 ;;
  darwin/x86_64)               OS=darwin; ARCH=amd64 ;;
  darwin/amd64)                ARCH=amd64 ;;
  linux/arm64|linux/aarch64)   OS=linux; ARCH=arm64 ;;
  linux/x86_64|linux/amd64)    OS=linux; ARCH=amd64 ;;
  mingw*/msys*/cygwin*)
    echo "windows: this installer is unix-only. Download the .zip from:" >&2
    echo "  https://github.com/$REPO/releases/latest" >&2
    exit 2 ;;
  *) echo "unsupported platform: $OS/$ARCH" >&2; exit 2 ;;
esac

# --- 2. resolve version ---
if [ -n "$PIN_VERSION" ]; then
  TAG="${PIN_VERSION#v}"
  TAG="v$TAG"
else
  echo "resolving latest release..."
  TAG="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
         | sed -nE 's/.*"tag_name"[[:space:]]*:[[:space:]]*"(v[^"]+)".*/\1/p' | head -1)"
  if [ -z "$TAG" ]; then
    echo "could not resolve latest release tag" >&2
    exit 1
  fi
fi
VERSION="${TAG#v}"
echo "installing $TAG ($OS/$ARCH)"

# --- 3. download archive + checksums ---
BASE="https://github.com/$REPO/releases/download/$TAG"
ASSET="adp_${VERSION}_${OS}_${ARCH}.tar.gz"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

echo "downloading ${ASSET}..."
curl -fsSL -o "$TMP/$ASSET"          "$BASE/$ASSET"
curl -fsSL -o "$TMP/checksums.txt"   "$BASE/adp_${VERSION}_checksums.txt"

# --- 4. verify SHA256 ---
EXPECTED="$(awk -v f="$ASSET" '$2==f {print $1}' "$TMP/checksums.txt")"
if [ -z "$EXPECTED" ]; then
  echo "no checksum entry for $ASSET in checksums.txt; aborting" >&2
  exit 3
fi
ACTUAL="$(sha256sum "$TMP/$ASSET" 2>/dev/null | awk '{print $1}' \
          || shasum -a 256 "$TMP/$ASSET" | awk '{print $1}')"
# lowercase both for comparison (shasum emits uppercase on some systems)
EXPECTED_L="$(printf '%s' "$EXPECTED" | tr '[:upper:]' '[:lower:]')"
ACTUAL_L="$(printf '%s' "$ACTUAL"   | tr '[:upper:]' '[:lower:]')"
if [ "$EXPECTED_L" != "$ACTUAL_L" ]; then
  echo "checksum mismatch for $ASSET" >&2
  echo "  have: $ACTUAL_L" >&2
  echo "  want: $EXPECTED_L" >&2
  exit 3
fi
echo "checksum ok"

# --- 5. extract + install ---
tar -xzf "$TMP/$ASSET" -C "$TMP" adp
mkdir -p "$INSTALL_DIR"
install -m 0755 "$TMP/adp" "$INSTALL_DIR/adp"

# --- 6. PATH hint ---
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo
    echo "note: $INSTALL_DIR is not on your PATH."
    echo "add this to your shell profile and restart the shell:"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    ;;
esac

echo
echo "installed → $INSTALL_DIR/adp ($TAG)"
echo "verify:    adp version"
echo "upgrade:   adp update"
