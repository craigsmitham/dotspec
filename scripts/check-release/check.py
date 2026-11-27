"""Verify the documentation and manual page shipped in each GoReleaser platform archive."""
import json
from pathlib import Path
import tarfile
import zipfile


def check(root):
    handbook = root / "docs/cli-handbook"
    pages = {path.relative_to(root).as_posix() for path in handbook.glob("*/**/*.md")}
    pages.update({"docs/cli-handbook/README.md", "docs/cli-handbook/topics.json"})
    pages.update(path.relative_to(root).as_posix() for path in (root / "docs/reference").rglob("*.md"))
    examples = {path.relative_to(root).as_posix() for path in (root / "examples/todo").rglob("*.spec.md")}
    examples.update(path.relative_to(root).as_posix() for path in (root / "examples/todo").rglob("*.quality-model.md"))
    dist = root / "dist"
    archives = sorted([*dist.glob("dotspec_*.tar.gz"), *dist.glob("dotspec_*.zip")])
    version = json.loads((dist / "metadata.json").read_text())["version"]
    expected = {f"dotspec_{version}_{os}_{arch}.{ext}"
                for os, ext in (("linux", "tar.gz"), ("darwin", "tar.gz"), ("windows", "zip"))
                for arch in ("amd64", "arm64")}
    if {archive.name for archive in archives} != expected:
        raise ValueError(f"Expected six platform archives for {version} in {dist}")
    for archive in archives:
        if archive.suffix == ".zip":
            with zipfile.ZipFile(archive) as bundle:
                names = set(bundle.namelist())
        else:
            with tarfile.open(archive) as bundle:
                names = set(bundle.getnames())
        exe = "dotspec.exe" if archive.suffix == ".zip" else "dotspec"
        missing = (pages | examples | {exe, "LICENSE", "README.md", "man/man1/dotspec.1"}) - names
        extra = {name for name in names if name.startswith("docs/")} - pages
        extra.update(names & {"CONTRIBUTING.md", "model/README.md", "examples/todo/COVERAGE.md"})
        extra.update(name for name in names if Path(name).name == "AGENTS.md")
        if missing or extra:
            raise ValueError(f"{archive.name}: missing {sorted(missing)}, unexpected {sorted(extra)}")
    print(f"Verified documentation in {len(archives)} platform archives.")


if __name__ == "__main__":
    check(Path(__file__).resolve().parents[2])
