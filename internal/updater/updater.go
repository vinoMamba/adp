// Package updater implements self-update for the adp CLI by reading the
// project's GitHub releases (produced by goreleaser) and atomically
// replacing the running executable.
//
// Only the standard library is used. Network access goes through the
// caller-provided http.Client (see DefaultClient) so tests can swap it.
package updater

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Repo is the GitHub owner/repo to query. Override at link time or in tests
// to point at a fork; defaults to buildinfo.Repo.
var Repo = "vinoMamba/adp"

// DefaultClient is used when Options.Client is nil. The 30s timeout covers
// large binary downloads on slow links because Download uses its own
// per-stage context; this client is only for the small API JSON calls.
var DefaultClient = &http.Client{Timeout: 30 * time.Second}

// MaxDownloadBytes caps any single asset fetch. See fetchBytes.
var MaxDownloadBytes int64 = 256 << 20

// Errors returned by this package. Compare with errors.Is.
var (
	// ErrNotFound is returned when the requested release or a matching
	// asset for the current GOOS/GOARCH cannot be located.
	ErrNotFound = errors.New("updater: release or asset not found")
	// ErrUpToDate is returned by Update when the latest release is not
	// newer than the running version (and Force is false).
	ErrUpToDate = errors.New("updater: already up to date")
)

// Asset is a single downloadable file attached to a release.
type Asset struct {
	Name string // e.g. "adp_0.2.0_darwin_arm64.tar.gz"
	URL  string // browser_download_url
	Size int64
}

// Release describes a GitHub release, filtered to the fields we need.
type Release struct {
	Tag      string    // "v0.2.0"
	Assets   []Asset   // downloadable files
	HTMLURL  string    // human-readable release page
	Body     string    // release notes markdown
	Published time.Time // publication timestamp
}

// assetName is the binary we look for inside an archive.
func assetBinaryName(goos string) string {
	if goos == "windows" {
		return BinaryNameVar + ".exe"
	}
	return BinaryNameVar
}

// BinaryNameVar is the executable base name expected inside the archive.
// Defaults to "adp"; exposed so tests can override.
var BinaryNameVar = "adp"

// Latest fetches metadata for the latest release (when version is empty) or
// the release tagged `version` (e.g. "v0.2.0" or "0.2.0"; the leading "v" is
// optional). The leading "v" is normalized on the returned Tag.
func Latest(ctx context.Context, client *http.Client, repo, version string) (*Release, error) {
	if client == nil {
		client = DefaultClient
	}
	if repo == "" {
		repo = Repo
	}
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	if version != "" {
		// GitHub accepts tags without the leading "v", but the canonical
		// tag form on this repo is "vX.Y.Z". Normalize so a user typing
		// `--version 0.2.0` still hits the right tag.
		tag := strings.TrimPrefix(version, "v")
		endpoint = fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/v%s", repo, tag)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	// UA helps GitHub rate-limit identification; harmless if absent.
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch release: %w", err)
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusNotFound:
		return nil, fmt.Errorf("%w: %s", ErrNotFound, endpoint)
	case resp.StatusCode != http.StatusOK:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("github api %s: %s", resp.Status, bytes.TrimSpace(body))
	}

	var raw struct {
		TagName     string    `json:"tag_name"`
		HTMLURL     string    `json:"html_url"`
		Body        string    `json:"body"`
		PublishedAt time.Time `json:"published_at"`
		Assets []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
			Size int64  `json:"size"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode release: %w", err)
	}
	if raw.TagName == "" {
		return nil, fmt.Errorf("%w: empty tag_name", ErrNotFound)
	}

	rel := &Release{
		Tag:       raw.TagName,
		HTMLURL:   raw.HTMLURL,
		Body:      raw.Body,
		Published: raw.PublishedAt,
	}
	for _, a := range raw.Assets {
		rel.Assets = append(rel.Assets, Asset{Name: a.Name, URL: a.URL, Size: a.Size})
	}
	return rel, nil
}

// NeedUpdate reports whether `latest` is a newer semver than `current`.
// Strings without a leading "v" and trailing pre-release/build metadata are
// tolerated. A missing/empty current is treated as "always update".
// Pre-release versions (e.g. v1.0.0-rc1) are considered lower than the same
// version without a pre-release, per semver.
func NeedUpdate(current, latest string) bool {
	if latest == "" {
		return false
	}
	if current == "" {
		return true
	}
	switch compareSemver(current, latest) {
	case -1:
		return true
	default:
		return false
	}
}

// MatchAsset finds the asset for the given goos/goarch in the release.
// It excludes checksums.txt and matches the "{goos}_{goarch}" substring,
// which is robust to goreleaser's archive naming variations.
func (r *Release) MatchAsset(goos, goarch string) (Asset, error) {
	want := "_" + goos + "_" + goarch
	var fallback Asset
	found := false
	for _, a := range r.Assets {
		name := strings.ToLower(a.Name)
		if strings.HasSuffix(name, "_checksums.txt") {
			continue
		}
		if !strings.Contains(name, want) {
			continue
		}
		// Prefer a suffix match ("..._darwin_arm64.tar.gz") over a
		// substring-only match, but accept the first plausible hit.
		if strings.HasSuffix(name, want+".tar.gz") || strings.HasSuffix(name, want+".zip") {
			return a, nil
		}
		if !found {
			fallback = a
			found = true
		}
	}
	if found {
		return fallback, nil
	}
	return Asset{}, fmt.Errorf("%w: no asset for %s/%s in %s", ErrNotFound, goos, goarch, r.Tag)
}

// checksumFor parses goreleaser's checksums.txt and returns the hex digest
// for the named asset. Returns "" (no error) when the file has no entry —
// callers decide whether to enforce.
func checksumFor(checksums []byte, asset string) string {
	for _, line := range bytes.Split(checksums, []byte("\n")) {
		fields := bytes.Fields(line)
		if len(fields) != 2 {
			continue
		}
		if string(fields[1]) == asset {
			return strings.ToLower(strings.TrimSpace(string(fields[0])))
		}
	}
	return ""
}

// Download fetches the asset matching goos/goarch. When the release also
// ships a checksums.txt, the downloaded archive is verified against it; a
// mismatch aborts the update. Returns the extracted executable bytes.
func (r *Release) Download(ctx context.Context, client *http.Client, goos, goarch string) ([]byte, error) {
	if client == nil {
		client = DefaultClient
	}
	asset, err := r.MatchAsset(goos, goarch)
	if err != nil {
		return nil, err
	}
	archive, err := fetchBytes(ctx, client, asset.URL)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", asset.Name, err)
	}

	// Verify SHA256 if checksums.txt is published alongside.
	var checksums []byte
	for _, a := range r.Assets {
		if strings.HasSuffix(strings.ToLower(a.Name), "_checksums.txt") {
			checksums, _ = fetchBytes(ctx, client, a.URL)
			break
		}
	}
	if want := checksumFor(checksums, asset.Name); want != "" {
		got := sha256Hex(archive)
		if !strings.EqualFold(got, want) {
			return nil, fmt.Errorf("checksum mismatch for %s: have %s, want %s", asset.Name, got, want)
		}
	}

	bin, err := extractBinary(bytes.NewReader(archive), asset.Name, assetBinaryName(goos))
	if err != nil {
		return nil, fmt.Errorf("extract %s: %w", asset.Name, err)
	}
	return bin, nil
}

// fetchBytes downloads a URL into memory with a sane size cap.
func fetchBytes(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %s: %s", url, resp.Status)
	}
	// MaxDownloadBytes caps any single asset fetch to block a malicious or
	// misconfigured asset from exhausting memory. Exported so tests can
	// lower it; 256 MiB is comfortably above our ~22 MiB release.
	return io.ReadAll(io.LimitReader(resp.Body, MaxDownloadBytes))
}

// extractBinary reads an archive (tar.gz or zip) and returns the bytes of
// the first entry whose base name matches `binary` (case-insensitive).
func extractBinary(r io.ReaderAt, name, binary string) ([]byte, error) {
	binLower := strings.ToLower(binary)
	switch {
	case strings.HasSuffix(strings.ToLower(name), ".zip"):
		return extractZip(r, binLower)
	case strings.HasSuffix(strings.ToLower(name), ".tar.gz"),
		strings.HasSuffix(strings.ToLower(name), ".tgz"):
		return extractTarGz(r, binLower)
	default:
		return nil, fmt.Errorf("unsupported archive format: %s", name)
	}
}

// extractZip assumes `r` is *bytes.Reader-compatible (the caller passes the
// full archive). archive/zip needs a ReaderAt and size.
func extractZip(r io.ReaderAt, binLower string) ([]byte, error) {
	// We don't know the size without re-reading; the caller passes a
	// *bytes.Reader via bytes.NewReader(archive). Use its Size() if available.
	size, err := readerSize(r)
	if err != nil {
		return nil, err
	}
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	for _, f := range zr.File {
		if strings.ToLower(filepath.Base(f.Name)) != binLower {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		return io.ReadAll(rc)
	}
	return nil, fmt.Errorf("%w: %s not in archive", ErrNotFound, binLower)
}

func extractTarGz(r io.ReaderAt, binLower string) ([]byte, error) {
	size, err := readerSize(r)
	if err != nil {
		return nil, err
	}
	sr := io.NewSectionReader(r, 0, size)
	gz, err := gzip.NewReader(sr)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		if strings.ToLower(filepath.Base(hdr.Name)) != binLower {
			continue
		}
		return io.ReadAll(tr)
	}
	return nil, fmt.Errorf("%w: %s not in archive", ErrNotFound, binLower)
}

// readerSize extracts the underlying size for the common case where the
// ReaderAt is a *bytes.Reader (which is what Release.Download passes).
func readerSize(r io.ReaderAt) (int64, error) {
	type sizer interface{ Size() int64 }
	if s, ok := r.(sizer); ok {
		return s.Size(), nil
	}
	// bytes.Reader implements Len() but not Size(); try Len.
	type lener interface{ Len() int }
	if l, ok := r.(lener); ok {
		return int64(l.Len()), nil
	}
	return 0, errors.New("updater: archive source must expose Size() or Len() (use bytes.NewReader)")
}

// Apply replaces the executable at `target` with `body`.
//
// On Unix, the replacement is atomic: write to a sibling temp file, chmod,
// then rename over the original. On Windows, the running binary cannot be
// overwritten in place, so the existing file is renamed to `*.old` first;
// the next launch best-effort deletes the stale copy.
//
// The original file's permission bits are preserved.
func Apply(target string, body []byte) error {
	if target == "" {
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("locate executable: %w", err)
		}
		target = exe
	}

	info, err := os.Stat(target)
	var mode os.FileMode = 0o755
	if err == nil {
		mode = info.Mode()
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat %s: %w", target, err)
	}

	dir := filepath.Dir(target)

	if runtime.GOOS == "windows" {
		return applyWindows(target, dir, body, mode)
	}
	return applyPosix(target, dir, body, mode)
}

func applyPosix(target, dir string, body []byte, mode os.FileMode) error {
	tmp, err := os.CreateTemp(dir, ".adp-update-*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpName := tmp.Name()
	cleanup := func() {
		tmp.Close()
		os.Remove(tmpName)
	}

	if _, err := tmp.Write(body); err != nil {
		cleanup()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Chmod(mode); err != nil {
		cleanup()
		return fmt.Errorf("chmod temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Rename(tmpName, target); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("replace %s: %w", target, err)
	}
	return nil
}

func applyWindows(target, dir string, body []byte, mode os.FileMode) error {
	// Move the locked, running binary aside, then write the new one to the
	// canonical path. The .old file is cleaned up on next launch via
	// CleanupStaleWindowsBackup.
	old := target + ".old"
	os.Remove(old) // remove any leftover from a previous run
	if err := os.Rename(target, old); err != nil {
		return fmt.Errorf("stage old binary: %w", err)
	}
	if err := os.WriteFile(target, body, mode); err != nil {
		// Roll back so the user isn't left without a working adp.
		_ = os.Rename(old, target)
		return fmt.Errorf("write new binary: %w", err)
	}
	// Best-effort; Windows may still hold the handle.
	_ = os.Remove(old)
	return nil
}

// CleanupStaleWindowsBackup removes any `adp.exe.old` left by a previous
// self-update. Safe to call on every launch; no-op on non-Windows.
func CleanupStaleWindowsBackup() {
	if runtime.GOOS != "windows" {
		return
	}
	exe, err := os.Executable()
	if err != nil {
		return
	}
	os.Remove(exe + ".old")
}

// sha256Hex returns the lowercase hex SHA256 of `data`.
func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// --- semver ---

var semverRe = regexp.MustCompile(`^v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:[-+](.*))?$`)

// compareSemver returns -1, 0, or +1 like bytes.Compare. Inputs may carry a
// leading "v". A non-parseable version compares as less than a parseable one.
func compareSemver(a, b string) int {
	av, okA := parseSemver(a)
	bv, okB := parseSemver(b)
	switch {
	case okA && !okB:
		return 1
	case !okA && okB:
		return -1
	case !okA && !okB:
		return strings.Compare(a, b)
	}
	if c := cmpInt(av.major, bv.major); c != 0 {
		return c
	}
	if c := cmpInt(av.minor, bv.minor); c != 0 {
		return c
	}
	if c := cmpInt(av.patch, bv.patch); c != 0 {
		return c
	}
	// Per semver, a version with a pre-release is LOWER than the same
	// version without. "1.0.0-rc1" < "1.0.0".
	if av.pre == "" && bv.pre != "" {
		return 1
	}
	if av.pre != "" && bv.pre == "" {
		return -1
	}
	return strings.Compare(av.pre, bv.pre)
}

type semver struct {
	major, minor, patch int
	pre                 string
}

func parseSemver(s string) (semver, bool) {
	s = strings.TrimSpace(s)
	m := semverRe.FindStringSubmatch(s)
	if m == nil {
		return semver{}, false
	}
	out := semver{
		major: atoiOr(m[1], 0),
		minor: atoiOr(m[2], 0),
		patch: atoiOr(m[3], 0),
		pre:   m[4],
	}
	return out, true
}

func atoiOr(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
