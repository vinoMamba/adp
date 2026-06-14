// Package buildinfo holds version metadata injected at link time.
package buildinfo

import "runtime/debug"

// Version, Commit, Date are overridden via -ldflags -X at release build time.
var (
	Version = ""
	Commit  = ""
	Date    = ""
)

// Repo is the GitHub owner/repo used by the self-updater to find releases.
const Repo = "vinoMamba/adp"

// BinaryName is the published executable name (without extension); used by
// the updater to locate the binary inside the release archive.
const BinaryName = "adp"

// Resolve returns the effective version: the injected Version if set, else the
// module version from the build info (so `go install`ed binaries self-describe).
func Resolve() string {
	if Version != "" {
		return Version
	}
	if bi, ok := debug.ReadBuildInfo(); ok {
		if bi.Main.Version != "" && bi.Main.Version != "(devel)" {
			return bi.Main.Version
		}
	}
	return ""
}

// String returns a human-readable version string.
func String() string {
	v := Resolve()
	if v == "" {
		v = "dev"
	}
	s := v
	if Commit != "" {
		s += " (" + Commit + ")"
	}
	if Date != "" {
		s += " " + Date
	}
	return s
}
