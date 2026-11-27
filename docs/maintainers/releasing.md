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

Push a SemVer tag on a commit that passed the hosted CI matrix. Use
`v0.1.0-alpha.1`, `v0.1.0-alpha.2`, and so on for alpha releases.
The release workflow builds a GitHub draft, packs the six checksummed executables
into one npm tarball, then installs that tarball on Linux, macOS, and Windows.
Each runner checks the version, opens the embedded handbook, and initializes,
builds, and checks a project. A passing draft becomes a public GitHub release.
Alpha releases carry the GitHub prerelease flag. Follow the tag run and verify
the six platform archives, knowledge ZIP, checksums, and npm tarball.
[Release workflow](../../.github/workflows/release.yml)

The Homebrew tap polls published releases every six hours. Run
`gh workflow run update.yml -R dotspec-md/homebrew-tap` to refresh immediately.
The tap installs and tests the formula before committing its new tag and revision.
[TAP workflow](https://github.com/dotspec-md/homebrew-tap/blob/main/.github/workflows/update.yml)

For the first npm release, download the verified `dotspec-<version>.tgz` from
the GitHub release. With the owning npm account signed in, run
`npm publish ./dotspec-<version>.tgz --access public --tag alpha` and complete
npm's authentication prompt. Configure npm trusted publishing for the
`dotspec` package, GitHub repository `dotspec-md/dotspec`, and workflow
`release.yml` with publish permission. Set the repository Actions variable
`NPM_TRUSTED_PUBLISHING=true` to enable subsequent tag publishes. The workflow
uses GitHub OIDC and maps alpha versions to the npm `alpha` dist-tag.
[Trusted publishing](https://docs.npmjs.com/trusted-publishers/)

Retry a failed smoke test or draft job from the tag run. Once public, keep that
tag and its assets intact. If only npm publication failed, dispatch
`gh workflow run release.yml -f tag=vX.Y.Z` for the existing release. Check
`npm view dotspec@<version> dist.integrity` and `npx --yes dotspec@<version> --version`
after publishing. Issue a new version when a distributed artifact needs a fix.
