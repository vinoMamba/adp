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
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCompareSemver(t *testing.T) {
	cases := []struct {
		name string
		a, b string
		want int
	}{
		{"equal", "v0.1.0", "0.1.0", 0},
		{"major", "v1.0.0", "v0.9.9", 1},
		{"minor", "0.1.10", "0.1.9", 1},
		{"patch", "v0.1.0", "v0.1.1", -1},
		{"missing_minor", "1", "1.0", 0},
		{"pre_lower", "1.0.0-rc1", "1.0.0", -1},
		{"pre_higher", "1.0.0", "1.0.0-rc1", 1},
		{"pre_compare", "1.0.0-rc1", "1.0.0-rc2", -1},
		{"unparseable_a", "garbage", "1.0.0", -1},
		{"unparseable_b", "1.0.0", "garbage", 1},
		{"both_unparseable", "abc", "abc", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := compareSemver(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("compareSemver(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestNeedUpdate(t *testing.T) {
	cases := []struct {
		current, latest string
		want            bool
	}{
		{"v0.1.0", "v0.2.0", true},
		{"v0.2.0", "v0.2.0", false},
		{"v0.3.0", "v0.2.0", false},
		{"", "v0.1.0", true},
		{"v0.1.0", "", false},
		{"v1.0.0-rc1", "v1.0.0", true},
	}
	for _, tc := range cases {
		name := tc.current + "_to_" + tc.latest
		t.Run(name, func(t *testing.T) {
			if got := NeedUpdate(tc.current, tc.latest); got != tc.want {
				t.Errorf("NeedUpdate(%q, %q) = %v, want %v", tc.current, tc.latest, got, tc.want)
			}
		})
	}
}

func TestMatchAsset(t *testing.T) {
	rel := &Release{
		Tag: "v0.2.0",
		Assets: []Asset{
			{Name: "adp_0.2.0_checksums.txt"},
			{Name: "adp_0.2.0_darwin_amd64.tar.gz"},
			{Name: "adp_0.2.0_darwin_arm64.tar.gz"},
			{Name: "adp_0.2.0_linux_amd64.tar.gz"},
			{Name: "adp_0.2.0_windows_amd64.zip"},
		},
	}
	cases := []struct {
		goos, goarch string
		wantName     string
		wantErr      bool
	}{
		{"darwin", "arm64", "adp_0.2.0_darwin_arm64.tar.gz", false},
		{"darwin", "amd64", "adp_0.2.0_darwin_amd64.tar.gz", false},
		{"linux", "amd64", "adp_0.2.0_linux_amd64.tar.gz", false},
		{"windows", "amd64", "adp_0.2.0_windows_amd64.zip", false},
		{"linux", "arm64", "", true}, // not shipped
	}
	for _, tc := range cases {
		name := tc.goos + "_" + tc.goarch
		t.Run(name, func(t *testing.T) {
			got, err := rel.MatchAsset(tc.goos, tc.goarch)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %+v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Name != tc.wantName {
				t.Errorf("got %q, want %q", got.Name, tc.wantName)
			}
		})
	}
}

func TestChecksumFor(t *testing.T) {
	checksums := []byte(strings.Join([]string{
		"abc123  adp_0.2.0_darwin_arm64.tar.gz",
		"def456  adp_0.2.0_linux_amd64.tar.gz",
		"", // trailing newline produces an empty final line
	}, "\n"))
	cases := []struct {
		asset string
		want  string
	}{
		{"adp_0.2.0_darwin_arm64.tar.gz", "abc123"},
		{"adp_0.2.0_linux_amd64.tar.gz", "def456"},
		{"missing.tar.gz", ""},
	}
	for _, tc := range cases {
		t.Run(tc.asset, func(t *testing.T) {
			if got := checksumFor(checksums, tc.asset); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// makeTarGz returns a tar.gz whose single regular file is named `name` with
// the given body. Used to exercise extractTarGz and the checksum path.
func makeTarGz(t *testing.T, name string, body []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{
		Name: name, Mode: 0o755, Size: int64(len(body)), Typeflag: tar.TypeReg,
	}); err != nil {
		t.Fatalf("write header: %v", err)
	}
	if _, err := tw.Write(body); err != nil {
		t.Fatalf("write body: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}
	return buf.Bytes()
}

// makeZip returns a zip whose single file is named `name` with the given body.
func makeZip(t *testing.T, name string, body []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create(name)
	if err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	if _, err := w.Write(body); err != nil {
		t.Fatalf("write zip body: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	return buf.Bytes()
}

func TestExtractBinaryTarGz(t *testing.T) {
	want := []byte("binary-bytes")
	archive := makeTarGz(t, "adp", want)
	got, err := extractBinary(bytes.NewReader(archive), "adp_0.2.0_linux_amd64.tar.gz", "adp")
	if err != nil {
		t.Fatalf("extractBinary: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("extracted bytes mismatch: have %q, want %q", got, want)
	}
}

func TestExtractBinaryZip(t *testing.T) {
	want := []byte("windows-bytes")
	archive := makeZip(t, "adp.exe", want)
	got, err := extractBinary(bytes.NewReader(archive), "adp_0.2.0_windows_amd64.zip", "adp.exe")
	if err != nil {
		t.Fatalf("extractBinary: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("extracted bytes mismatch: have %q, want %q", got, want)
	}
}

func TestExtractBinaryNestedName(t *testing.T) {
	// goreleaser sometimes nests archives; the lookup must use Base().
	want := []byte("nested")
	archive := makeTarGz(t, "adp_0.2.0_darwin_arm64/adp", want)
	got, err := extractBinary(bytes.NewReader(archive), "adp_0.2.0_darwin_arm64.tar.gz", "adp")
	if err != nil {
		t.Fatalf("extractBinary: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("have %q, want %q", got, want)
	}
}

// A fake GitHub + download server lets us exercise Download end-to-end,
// including SHA256 verification via a published checksums.txt.
func TestDownloadVerifiesChecksum(t *testing.T) {
	wantBin := []byte("the-real-binary")
	archiveName := "adp_0.2.0_darwin_arm64.tar.gz"
	archive := makeTarGz(t, "adp", wantBin)
	sum := sha256.Sum256(archive)
	checksums := []byte(fmt.Sprintf("%s  %s\n", hex.EncodeToString(sum[:]), archiveName))

	mux := http.NewServeMux()
	mux.HandleFunc("/archive", func(w http.ResponseWriter, r *http.Request) {
		w.Write(archive)
	})
	mux.HandleFunc("/checksums.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Write(checksums)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	rel := &Release{
		Tag: "v0.2.0",
		Assets: []Asset{
			{Name: archiveName, URL: srv.URL + "/archive", Size: int64(len(archive))},
			{Name: "adp_0.2.0_checksums.txt", URL: srv.URL + "/checksums.txt"},
		},
	}

	got, err := rel.Download(context.Background(), srv.Client(), "darwin", "arm64")
	if err != nil {
		t.Fatalf("Download: %v", err)
	}
	if !bytes.Equal(got, wantBin) {
		t.Errorf("have %q, want %q", got, wantBin)
	}
}

func TestDownloadRejectsCorruptArchive(t *testing.T) {
	archiveName := "adp_0.2.0_darwin_arm64.tar.gz"
	archive := makeTarGz(t, "adp", []byte("good"))
	// Publish a checksum for a DIFFERENT archive — verification must fail.
	checksums := []byte("0000000000000000000000000000000000000000000000000000000000000000  " + archiveName + "\n")

	mux := http.NewServeMux()
	mux.HandleFunc("/archive", func(w http.ResponseWriter, r *http.Request) { w.Write(archive) })
	mux.HandleFunc("/checksums.txt", func(w http.ResponseWriter, r *http.Request) { w.Write(checksums) })
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	rel := &Release{
		Tag: "v0.2.0",
		Assets: []Asset{
			{Name: archiveName, URL: srv.URL + "/archive"},
			{Name: "adp_0.2.0_checksums.txt", URL: srv.URL + "/checksums.txt"},
		},
	}
	if _, err := rel.Download(context.Background(), srv.Client(), "darwin", "arm64"); err == nil {
		t.Fatal("expected checksum mismatch error, got nil")
	}
}

func TestApplyPosixAtomicReplace(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("posix-only behavior")
	}
	dir := t.TempDir()
	target := filepath.Join(dir, "adp")
	original := []byte("v1")
	if err := os.WriteFile(target, original, 0o755); err != nil {
		t.Fatalf("seed: %v", err)
	}

	updated := []byte("v2")
	if err := Apply(target, updated); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !bytes.Equal(got, updated) {
		t.Errorf("have %q, want %q", got, updated)
	}
	// No leftover temp files.
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".adp-update-") {
			t.Errorf("leftover temp file: %s", e.Name())
		}
	}
}

func TestApplyPreservesMode(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "adp")
	if err := os.WriteFile(target, []byte("x"), 0o700); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := Apply(target, []byte("y")); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o700 {
		t.Errorf("mode = %o, want 0700", got)
	}
}

// TestLatestUsesCustomClient exercises Latest end-to-end against a stub
// serving GitHub-shaped JSON.
func TestLatestUsesCustomClient(t *testing.T) {
	payload := map[string]any{
		"tag_name":     "v0.3.0",
		"html_url":     "https://example/release",
		"body":         "release notes",
		"published_at": "2026-06-01T00:00:00Z",
		"assets": []map[string]any{
			{"name": "adp_0.3.0_darwin_arm64.tar.gz", "browser_download_url": "https://example/a", "size": 42},
		},
	}
	body, _ := json.Marshal(payload)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mirror the latest-release endpoint regardless of repo path.
		if !strings.HasSuffix(r.URL.Path, "/releases/latest") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}))
	t.Cleanup(srv.Close)

	// Build a client that rewrites api.github.com → our stub.
	client := &http.Client{Transport: rewriteTransport{base: srv.URL}}
	rel, err := Latest(context.Background(), client, "owner/repo", "")
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if rel.Tag != "v0.3.0" {
		t.Errorf("Tag = %q, want v0.3.0", rel.Tag)
	}
	if len(rel.Assets) != 1 || rel.Assets[0].Name != "adp_0.3.0_darwin_arm64.tar.gz" {
		t.Errorf("Assets = %+v", rel.Assets)
	}
}

func TestLatestWrapsNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	t.Cleanup(srv.Close)

	client := &http.Client{Transport: rewriteTransport{base: srv.URL}}
	_, err := Latest(context.Background(), client, "owner/repo", "v9.9.9")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Latest error = %v, want ErrNotFound", err)
	}
}

// rewriteTransport rewrites every request URL's scheme+host to `base`, so a
// single httptest.Server can impersonate api.github.com and the download host.
type rewriteTransport struct{ base string }

func (t rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	stub := strings.TrimRight(t.base, "/")
	clone := req.Clone(req.Context())
	clone.URL.Scheme = strings.SplitN(stub, "://", 2)[0]
	clone.URL.Host = strings.SplitN(stub, "://", 2)[1]
	return http.DefaultTransport.RoundTrip(clone)
}

func TestFetchBytesSizeCap(t *testing.T) {
	// Lower the cap so the test doesn't allocate 256 MiB.
	original := MaxDownloadBytes
	MaxDownloadBytes = 1024
	t.Cleanup(func() { MaxDownloadBytes = original })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, endlessReader{})
	}))
	t.Cleanup(srv.Close)

	got, err := fetchBytes(context.Background(), srv.Client(), srv.URL)
	if err != nil {
		t.Fatalf("fetchBytes: %v", err)
	}
	if int64(len(got)) != MaxDownloadBytes {
		t.Errorf("len = %d, want %d", len(got), MaxDownloadBytes)
	}
}

type endlessReader struct{}

func (endlessReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 'x'
	}
	return len(p), nil
}
