"""Verify the documentation shipped in each GoReleaser platform archive."""
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
    if not archives:
        raise ValueError(f"No platform archives in {dist}")
    for archive in archives:
        if archive.suffix == ".zip":
            with zipfile.ZipFile(archive) as bundle:
                names = set(bundle.namelist())
        else:
            with tarfile.open(archive) as bundle:
                names = set(bundle.getnames())
        missing = (pages | examples | {"README.md"}) - names
        extra = {name for name in names if name.startswith("docs/")} - pages
        extra.update(names & {"CONTRIBUTING.md", "model/README.md", "examples/todo/COVERAGE.md"})
        extra.update(name for name in names if Path(name).name == "AGENTS.md")
        if missing or extra:
            raise ValueError(f"{archive.name}: missing {sorted(missing)}, unexpected {sorted(extra)}")
    print(f"Verified documentation in {len(archives)} platform archives.")


if __name__ == "__main__":
    check(Path(__file__).resolve().parents[2])
