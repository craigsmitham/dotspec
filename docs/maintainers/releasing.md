# Releasing dotspec

The release configuration builds standalone Go executables for Linux, macOS, and
Windows on amd64 and arm64 with CGO disabled. The binary embeds the handbook and examples. Each platform archive also
includes the Markdown guides for browsing and the manual page `man/man1/dotspec.1`, which a GoReleaser hook writes with the hidden `dotspec man` command. Releases attach a separate optional
`dotspec-knowledge.zip` for AXM users. GoReleaser generates `checksums.txt` and stamps the version and commit
reported by `dotspec --version`.

Work from the root of a clone with Git history and mise installed. Use a local snapshot to check what someone will download. The snapshot task refreshes `dist/`, so keep any earlier artifacts you need elsewhere. In a fresh checkout:

```sh
mise trust
mise install
mise run model-setup
mise run check
mise run fuzz
mise run snapshot
```

If a fuzz test fails, keep its reported input and resolve the failing assertion before tagging the revision.

The snapshot task checks each archive for the public handbook, worked specifications and manual page, and rejects repository-only documentation. From `dist/`, choose the archive matching the host OS and architecture. Compare its SHA-256 digest with the entry for that filename in `checksums.txt` (use `shasum -a 256 FILE` on macOS, `sha256sum FILE` on Linux, or `Get-FileHash FILE -Algorithm SHA256` in PowerShell). A mismatch warrants a fresh build before distribution. Extract the archive into an empty directory and run its executable with `--version` and `docs job`. Check that the build identifies the intended revision and the installed handbook opens. [GoReleaser checksums](https://goreleaser.com/customization/package/checksum/) · [Snapshot builds](https://goreleaser.com/customization/publish/snapshots/). Snapshot builds never publish and may be run from an uncommitted worktree.
Review and commit the final changes before creating a release tag.

A pushed `vX.Y.Z` tag triggers `.github/workflows/release.yml`. The workflow runs the
same full check task, then publishes the archives and checksums to the repository's
GitHub release using `GITHUB_TOKEN`. Repository Actions must be allowed to create releases. When you intend to publish, push the new version tag for the reviewed commit. Follow the run for that tag and open the resulting release. Check for six platform archives, `dotspec-knowledge.zip`, and `checksums.txt`. If publication fails, inspect the failed workflow step before deciding what to retry. [Release workflow](../../.github/workflows/release.yml) · [Archive configuration](../../.goreleaser.yaml).

Download an archive for the intended platform from GitHub Releases, verify it against
`checksums.txt`, extract it, and place `dotspec` (`dotspec.exe` on Windows) on PATH.
Remove that executable to uninstall. Users with Go can also install the repository
root package with `go install github.com/craigsmitham/dotspec@vX.Y.Z`, substituting
an available release tag.
