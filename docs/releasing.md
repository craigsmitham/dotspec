# Releasing dotspec

The release configuration builds standalone Go executables for Linux, macOS, and
Windows on amd64 and arm64 with CGO disabled. Each archive includes the CLI reference
and README. GoReleaser generates `checksums.txt` and stamps the version and commit
reported by `dotspec --version`.

Before preparing a release:

```sh
mise install
mise run check
mise run fuzz
mise run snapshot
```

Inspect `dist/`, verify the host executable's version and help, and verify archive
checksums. Snapshot builds never publish and may be run from an uncommitted worktree.
Review and commit the final changes before creating a release tag.

A pushed `vX.Y.Z` tag triggers `.github/workflows/release.yml`. The workflow runs the
same full check task, then publishes the archives and checksums to the repository's
GitHub release using `GITHUB_TOKEN`. Repository Actions must be allowed to create
releases. Creating or pushing a tag is a separate release action; adopting this
configuration does not publish a release.

Download an archive for the intended platform from GitHub Releases, verify it against
`checksums.txt`, extract it, and place `dotspec` (`dotspec.exe` on Windows) on PATH.
Remove that executable to uninstall. Users with Go can also install the repository
root package with `go install github.com/craigsmitham/dotspec@vX.Y.Z`, substituting
an available release tag.
